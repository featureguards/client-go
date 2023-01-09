package featureguards

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/featureguards/featureguards-go/v2/internal/client"
	"github.com/featureguards/featureguards-go/v2/internal/random"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	pb_ft "github.com/featureguards/featureguards-go/v2/proto/feature_toggle"
)

const (
	defaultDomain = "api.featureguards.com:443"
	dialTimeout   = 1 * time.Second
	apiTimeout    = 3 * time.Second
	errTimeout    = 1 * time.Minute
)

type featureToggles struct {
	client *client.Client

	// immutable
	defaults map[string]bool

	// atomic operations only
	clientVersion int64

	// errMu protects below
	errMu             sync.Mutex
	errDeadlineByName map[string]time.Time

	// mu protects below
	mu           sync.RWMutex
	ftByName     map[string]*pb_ft.FeatureToggle
	accessToken  string
	refreshToken string
}

// newFeatureToggles creates a new feature toggles client that fetches the latest state of feature toggles and
// returns. It also kicks off a go routine to sync changes to feature toggles in the background. The
// ctx passed in should be that controlling the entire lifetime of the client and not a per request
// context. This is because the client is expected to be long running.
func newFeatureToggles(ctx context.Context, options ...Options) (*featureToggles, error) {
	opts := &toggleOptions{}
	for _, opt := range options {
		if err := opt(opts); err != nil {
			return nil, err
		}
	}
	var clientOptions []client.Options
	if opts.dialOptions != nil {
		clientOptions = append(clientOptions, client.WithDialOptions(opts.dialOptions...))
	}
	if opts.apiKey != "" {
		clientOptions = append(clientOptions, client.WithApiKey(opts.apiKey))
	}
	if opts.domain == "" {
		opts.domain = defaultDomain
	}
	clientOptions = append(clientOptions, client.WithDomain(opts.domain))
	clientOptions = append(clientOptions, client.WithLogLevel(opts.logLevel))
	cl, err := client.New(ctx, clientOptions...)
	if err != nil {
		return nil, err
	}

	apiCtx, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()
	accessToken, refreshToken, err := cl.Authenticate(apiCtx)
	if err != nil {
		return nil, err
	}

	apiCtx, cancel = context.WithTimeout(ctx, apiTimeout)
	defer cancel()
	fetched, err := cl.Fetch(apiCtx, accessToken, int64(0))
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())

	toggles := &featureToggles{
		client:            cl,
		ftByName:          make(map[string]*pb_ft.FeatureToggle),
		errDeadlineByName: make(map[string]time.Time),
		accessToken:       accessToken,
		refreshToken:      refreshToken,
		defaults:          opts.defaults,
	}

	toggles.process(fetched.FeatureToggles, fetched.Version)

	if !opts.withoutListen {
		go toggles.listenLoop(ctx)
	}

	return toggles, nil
}

type ftOptions struct {
	attrs Attributes
}

// FeatureToggleOptions provides optional options to IsOn, such as attributes.
type FeatureToggleOptions func(o *ftOptions) error

// IsOn returns whether the feature toggle with the given name is on or not based on its settings and
// the passed in options, which include any attributes we're matching against.
func (ft *featureToggles) IsOn(name string, options ...FeatureToggleOptions) (bool, error) {
	defaults := ft.defaults[name]
	ft.mu.RLock()
	found, ok := ft.ftByName[name]
	ft.mu.RUnlock()
	if !ok {
		err := errors.Errorf("feature %s not found", name)
		ft.maybeLogError(name, err)
		return defaults, err
	}

	on, err := isOn(found, options...)
	if err != nil {
		ft.maybeLogError(name, err)
		return defaults, err
	}
	return on, nil
}

func isOn(ft *pb_ft.FeatureToggle, options ...FeatureToggleOptions) (bool, error) {
	opts := &ftOptions{}
	for _, opt := range options {
		if err := opt(opts); err != nil {
			return false, err
		}
	}

	if !ft.Enabled {
		return false, nil
	}
	if ft.DeletedAt.IsValid() {
		return false, errors.Errorf("feature toggle %s was deleted", ft.Name)
	}

	on := false
	switch ft.ToggleType {
	case pb_ft.FeatureToggle_ON_OFF:
		def := ft.GetOnOff()
		if def == nil || def.On == nil || def.Off == nil {
			return false, errors.Errorf("feature toggle %s is invalid", ft.Name)
		}
		if def.Off.Weight != 0 && def.Off.Weight != 100 || def.On.Weight != 0 && def.On.Weight != 100 ||
			def.On.Weight == def.Off.Weight {
			return false, errors.Errorf("invalid weights for feature toggle %s", ft.Name)
		}
		if def.On.Weight > 0 {
			on = true
		}
		// Now process allow list followed by disallow list
		if !on {
			matched, err := match(ft.Name, def.On.Matches, opts.attrs)
			if err != nil {
				return false, err
			}
			on = matched
		}
		if on {
			matched, err := match(ft.Name, def.Off.Matches, opts.attrs)
			if err != nil {
				return false, err
			}
			// This is a disallow list. Hence, should be negated.
			on = !matched
		}
	case pb_ft.FeatureToggle_PERCENTAGE:
		def := ft.GetPercentage()
		if def == nil || def.On == nil || def.Off == nil || def.Stickiness == nil {
			return false, errors.Errorf("feature toggle %s is invalid", ft.Name)
		}
		if def.On.Weight+def.Off.Weight != 100 || def.On.Weight < 0 || def.Off.Weight < 0 {
			return false, errors.Errorf("invalid weights for feature toggle %s", ft.Name)
		}

		switch def.Stickiness.StickinessType {
		case pb_ft.Stickiness_RANDOM:
			on = rand.Float32()*100 < def.On.Weight
		case pb_ft.Stickiness_KEYS:
			hash, err := hash(ft.Name, def.Stickiness.Keys, def.Salt, opts.attrs)
			if err != nil {
				return false, err
			}
			// Provides 0.000001 precision
			on = (hash % 1000000) < uint64(def.On.Weight*10000)
		}

		// Now process allow list followed by disallow list
		if !on {
			matched, err := match(ft.Name, def.On.Matches, opts.attrs)
			if err != nil {
				return false, err
			}
			on = matched
		}
		if on {
			matched, err := match(ft.Name, def.Off.Matches, opts.attrs)
			if err != nil {
				return false, err
			}
			// This is a disallow list. Hence, should be negated.
			on = !matched
		}

	}
	return on, nil
}

// maybeLogError avoids logging the same error over and over again to avoid performance issues.
func (ft *featureToggles) maybeLogError(name string, err error) {
	now := time.Now()
	ft.errMu.Lock()
	deadline, ok := ft.errDeadlineByName[name]
	if !ok || deadline.Before(now) {
		ft.errDeadlineByName[name] = now.Add(random.Jitter(errTimeout))
		ft.errMu.Unlock()
		log.Error(err)
	} else {
		ft.errMu.Unlock()
	}
}

func hash(name string, keys []*pb_ft.Key, salt string, attrs Attributes) (uint64, error) {
	if len(keys) == 0 {
		return 0, errors.Errorf("no attributes defined for feature toggle %s", name)
	}

	// process keys in order and pick the first one that exists
	for _, key := range keys {
		if key.Key == "" {
			return 0, errors.Errorf("no key name passed for feature toggle %s", name)
		}
		attr, ok := attrs[key.Key]
		if !ok {
			continue
		}
		v := salt
		// Below must be uniform across all languages
		switch key.KeyType {
		case pb_ft.Key_BOOLEAN:
			found, ok := attr.(bool)
			if !ok {
				return 0, errors.Errorf("expected boolean for key %s for feature toggle %s", key.Key, name)
			}
			if found {
				v += "true"
			} else {
				v += "false"
			}
		case pb_ft.Key_STRING:
			found, ok := attr.(string)
			if !ok {
				return 0, errors.Errorf("expected string for key %s for feature toggle %s", key.Key, name)
			}
			v += found
		case pb_ft.Key_FLOAT:
			found, ok := attr.(float32)
			if !ok {
				return 0, errors.Errorf("expected float for key %s for feature toggle %s", key.Key, name)
			}
			v += strconv.FormatFloat(float64(found), 'f', -1, 64)
		case pb_ft.Key_INT:
			found, ok := attr.(int64)
			if !ok {
				return 0, errors.Errorf("expected int/int64 for key %s for feature toggle %s", key.Key, name)
			}
			v += strconv.FormatInt(found, 10)
		case pb_ft.Key_DATE_TIME:
			found, ok := attr.(time.Time)
			if !ok {
				return 0, errors.Errorf("expected time.Time for key %s for feature toggle %s", key.Key, name)
			}
			v += strconv.FormatInt(found.UnixMilli(), 10)
		default:
			return 0, fmt.Errorf("unknown key type for %s for feature toggle %s", key.Key, name)
		}

		// Use md5 because it's also fast on the browser. Otherwise, would have opted in for murmur.
		return xxhash.Sum64([]byte(v)), nil
	}

	return 0, fmt.Errorf("no matching keys passed for feature toggle %s", name)
}

func match(name string, matches []*pb_ft.Match, attrs Attributes) (bool, error) {
	for _, m := range matches {
		if m.Key == nil || m.Key.Key == "" {
			// Bug
			return false, errors.Errorf("invalid match key for %s", name)
		}
		found := attrs[m.Key.Key]
		if found == nil {
			continue
		}
		switch m.Key.KeyType {
		case pb_ft.Key_BOOLEAN:
			v, ok := found.(bool)
			if !ok {
				return false, errors.Errorf("value passed for key %s is not boolean for feature toggle %s", m.Key.Key, name)
			}
			if m.GetBoolOp() == nil {
				return false, errors.Errorf("no boolean operation set for key %s and feature toggle %s", m.Key.Key, name)
			}
			value := m.GetBoolOp().Value
			if v == value {
				return true, nil
			}
		case pb_ft.Key_DATE_TIME:
			v, ok := found.(time.Time)
			if !ok || m.GetDateTimeOp() == nil {
				return false, errors.Errorf("value passed for key %s is not time.Time for feature toggle %s", m.Key.Key, name)
			}
			switch m.GetDateTimeOp().Op {
			case pb_ft.DateTimeOp_AFTER:
				on := v.After(m.GetDateTimeOp().Timestamp.AsTime())
				if on {
					return true, nil
				}
			case pb_ft.DateTimeOp_BEFORE:
				on := v.Before(m.GetDateTimeOp().Timestamp.AsTime())
				if on {
					return true, nil
				}
			default:
				// To support future operations, ignore it.
			}
		case pb_ft.Key_FLOAT:
			v, ok := found.(float32)
			if !ok || m.GetFloatOp() == nil {
				return false, errors.Errorf("value passed for key %s is not float for feature toggle %s", m.Key.Key, name)
			}
			values := m.GetFloatOp().Values
			switch op := m.GetFloatOp().Op; op {
			case pb_ft.FloatOp_IN:
				for _, value := range values {
					if v == value {
						return true, nil
					}
				}
			default:
				if len(values) != 1 {
					// Bug
					return false, errors.Errorf("invalid no. of values for key %s for feature toggle %s", m.Key.Key, name)
				}
				on := false
				switch op {
				case pb_ft.FloatOp_EQ:
					on = v == values[0]
				case pb_ft.FloatOp_GT:
					on = v > values[0]
				case pb_ft.FloatOp_GTE:
					on = v >= values[0]
				case pb_ft.FloatOp_LT:
					on = v < values[0]
				case pb_ft.FloatOp_LTE:
					on = v <= values[0]
				case pb_ft.FloatOp_NEQ:
					on = v != values[0]
				}
				if on {
					return true, nil
				}
			}
		case pb_ft.Key_INT:
			v, ok := found.(int64)
			if !ok || m.GetIntOp() == nil {
				return false, errors.Errorf("value passed for key %s is not int64 for feature toggle %s", m.Key.Key, name)
			}
			values := m.GetIntOp().Values
			switch op := m.GetIntOp().Op; op {
			case pb_ft.IntOp_IN:
				for _, value := range values {
					if v == value {
						return true, nil
					}
				}
			default:
				if len(values) != 1 {
					// Bug
					return false, errors.Errorf("invalid no. of values for key %s for feature toggle %s", m.Key.Key, name)
				}
				on := false
				switch op {
				case pb_ft.IntOp_EQ:
					on = v == values[0]
				case pb_ft.IntOp_GT:
					on = v > values[0]
				case pb_ft.IntOp_GTE:
					on = v >= values[0]
				case pb_ft.IntOp_LT:
					on = v < values[0]
				case pb_ft.IntOp_LTE:
					on = v <= values[0]
				case pb_ft.IntOp_NEQ:
					on = v != values[0]
				}
				if on {
					return true, nil
				}
			}
		case pb_ft.Key_STRING:
			v, ok := found.(string)
			if !ok || m.GetStringOp() == nil {
				return false, errors.Errorf("value passed for key %s is not string for feature toggle %s", m.Key.Key, name)
			}
			values := m.GetStringOp().Values
			switch op := m.GetStringOp().Op; op {
			case pb_ft.StringOp_IN:
				for _, value := range values {
					if value == v {
						return true, nil
					}
				}
			case pb_ft.StringOp_CONTAINS:
				if len(values) != 1 {
					// Bug
					return false, errors.Errorf("invalid no. of values for key %s for feature toggle %s", m.Key.Key, name)
				}
				if strings.Contains(v, values[0]) {
					return true, nil
				}
			case pb_ft.StringOp_EQ:
				if len(values) != 1 {
					// Bug
					return false, errors.Errorf("invalid no. of values for key %s for feature toggle %s", m.Key.Key, name)
				}
				if strings.Compare(v, values[0]) == 0 {
					return true, nil
				}
			default:
				// To support future operations, ignore it.
			}
		}
	}

	return false, nil
}
