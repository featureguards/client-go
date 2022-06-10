package featureguards

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/featureguards/featureguards-go/v1/internal/random"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	retryConnect = 5 * time.Second
)

var (
	ErrNoFeatureToggles error = errors.New("can't connect to feature guards")
)

type ResilientFeatureToggles struct {
	mu       sync.RWMutex
	ft       *featureToggles
	defaults map[string]bool
}

func New(ctx context.Context, options ...Options) *ResilientFeatureToggles {
	// extract the defaults
	opts := &toggleOptions{}
	for _, opt := range options {
		opt(opts)
	}
	cert := "../certs/prod.pem"
	if opts.testCerts {
		cert = "../certs/test.pem"
	}
	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		log.Errorln("Feature-guards couldn't be initialized. This should never happen.")
	}

	defaults := opts.defaults
	if defaults == nil {
		defaults = make(map[string]bool)
	}
	options = append(options, WithDialOptions(grpc.WithTransportCredentials(creds)))
	ft, err := newFeatureToggles(ctx, options...)
	r := &ResilientFeatureToggles{
		ft:       ft,
		defaults: defaults,
	}
	if err != nil {
		// Retry connecting in the background. Never block.
		go func() {
			for {
				select {
				case <-time.After(random.Jitter(retryConnect)):
					ft, err := newFeatureToggles(ctx, options...)
					if err == nil {
						r.mu.Lock()
						r.ft = ft
						r.mu.Unlock()
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	return r
}

// IsOn returns whether the feature toggle with the given name is on or not based on its settings and
// the passed in options, which include any attributes we're matching against.
func (r *ResilientFeatureToggles) IsOn(name string, options ...FeatureToggleOptions) (bool, error) {
	r.mu.RLock()
	ft := r.ft
	r.mu.RUnlock()
	if ft != nil {
		return ft.IsOn(name, options...)
	}
	found := r.defaults[name]
	return found, ErrNoFeatureToggles
}
