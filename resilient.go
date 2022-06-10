package featureguards

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/featureguards/featureguards-go/v2/certs"
	"github.com/featureguards/featureguards-go/v2/internal/random"
	"github.com/pkg/errors"
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

// New creates a new FeatureGuards client. The context passed in is expected to be long-running and
// controls the life-time of the client, usually the same lifetime as the binary.
// New dials in a separate go routine and will try to establish connection to FeatureGuards over time.
func New(ctx context.Context, options ...Options) *ResilientFeatureToggles {
	// extract the defaults
	opts := &toggleOptions{}
	for _, opt := range options {
		opt(opts)
	}
	if opts.domain == "" {
		opts.domain = defaultDomain
	}
	defaults := opts.defaults
	if defaults == nil {
		defaults = make(map[string]bool)
	}
	creds, err := tlsCreds(opts.testCerts, opts.domain)
	if err != nil {
		log.Error("Could not initialize feature-guards.")
	}
	options = append(options, WithDialOptions(grpc.WithTransportCredentials(creds)))
	ft, err := newFeatureToggles(ctx, options...)
	r := &ResilientFeatureToggles{
		ft:       ft,
		defaults: defaults,
	}
	if err != nil {
		log.Warnf("Could not initialize feature-guards due to %s. Will retry again.\n", err)
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
// the passed in options, which include any attributes FeatureGuards rules match against.
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

func tlsCreds(isTest bool, addr string) (credentials.TransportCredentials, error) {
	var cert []byte
	if isTest {
		cert = certs.TestCA
	} else {
		b, err := ioutil.ReadFile("/etc/ssl/cert.pem")
		if err != nil {
			return nil, errors.WithStack(err)
		}
		cert = b
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(cert) {
		return nil, errors.WithStack(errors.New("could not append cert"))
	}
	splitted := strings.Split(addr, ":")
	return credentials.NewTLS(&tls.Config{ServerName: splitted[0], RootCAs: cp}), nil
}
