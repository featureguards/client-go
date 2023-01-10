package featureguards

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/featureguards/featureguards-go/v2/internal/client"
	"github.com/featureguards/featureguards-go/v2/internal/random"
	log "github.com/sirupsen/logrus"
)

const (
	defaultDomain = "api.featureguards.com:443"
	dialTimeout   = 1 * time.Second
	apiTimeout    = 3 * time.Second
	errTimeout    = 1 * time.Minute
)

type clientWrapper struct {
	client *client.Client

	// immutable
	ft *featureToggles

	// atomic operations only
	ftVersion int64
	dsVersion int64

	// errMu protects below
	errMu             sync.Mutex
	errDeadlineByName map[string]time.Time

	// mu protects below
	mu           sync.RWMutex
	accessToken  string
	refreshToken string
}

// newClientWrapper creates a new feature toggles client that fetches the latest state of feature toggles and
// returns. It also kicks off a go routine to sync changes to feature toggles in the background. The
// ctx passed in should be that controlling the entire lifetime of the client and not a per request
// context. This is because the client is expected to be long running.
func newClientWrapper(ctx context.Context, options ...Options) (*clientWrapper, error) {
	rand.Seed(time.Now().UnixNano())
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
	fetched, err := cl.Fetch(apiCtx, accessToken, int64(0), int64(0))
	if err != nil {
		return nil, err
	}

	ft, err := newFeatureToggles(ctx, fetched.FeatureToggles, fetched.Version, options...)
	if err != nil {
		return nil, err
	}
	client := &clientWrapper{
		client:            cl,
		errDeadlineByName: make(map[string]time.Time),
		accessToken:       accessToken,
		refreshToken:      refreshToken,
		ft:                ft,
	}

	ft.process(fetched.FeatureToggles, fetched.Version)

	if !opts.withoutListen {
		go client.listenLoop(ctx)
	}

	return client, nil
}

// IsOn returns whether the feature toggle with the given name is on or not based on its settings and
// the passed in options, which include any attributes we're matching against.
func (cw *clientWrapper) IsOn(name string, options ...FeatureToggleOptions) (bool, error) {
	on, err := cw.ft.IsOn(name, options...)
	if err != nil {
		cw.maybeLogError(name, err)
	}
	return on, err
}

// maybeLogError avoids logging the same error over and over again to avoid performance issues.
func (cw *clientWrapper) maybeLogError(name string, err error) {
	now := time.Now()
	cw.errMu.Lock()
	deadline, ok := cw.errDeadlineByName[name]
	if !ok || deadline.Before(now) {
		cw.errDeadlineByName[name] = now.Add(random.Jitter(errTimeout))
		cw.errMu.Unlock()
		log.Error(err)
	} else {
		cw.errMu.Unlock()
	}
}
