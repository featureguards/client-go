package client

import (
	"context"
	"fmt"
	"time"

	"github.com/featureguards/featureguards-go/v2/internal/logger"
	"github.com/pkg/errors"

	pb_auth "github.com/featureguards/featureguards-go/v2/proto/auth"
	pb_toggles "github.com/featureguards/featureguards-go/v2/proto/toggles"

	"google.golang.org/grpc"
)

const (
	dialTimeout = 1 * time.Second
	apiTimeout  = 3 * time.Second
)

type clientOptions struct {
	domain      string
	logLevel    logger.LogLevel
	dialOptions []grpc.DialOption
	apiKey      string
}

type Options func(o *clientOptions) error

// For internal testing purposes mostly.
func WithDomain(domain string) Options {
	return func(o *clientOptions) error {
		o.domain = domain
		return nil
	}
}

func WithLogLevel(level logger.LogLevel) Options {
	return func(o *clientOptions) error {
		o.logLevel = level
		return nil
	}
}

func WithDialOptions(options ...grpc.DialOption) Options {
	return func(o *clientOptions) error {
		o.dialOptions = options
		return nil
	}
}

func WithApiKey(key string) Options {
	return func(o *clientOptions) error {
		o.apiKey = key
		return nil
	}
}

type Client struct {
	authClient    pb_auth.AuthClient
	togglesClient pb_toggles.TogglesClient
	apiKey        string
}

func New(ctx context.Context, options ...Options) (*Client, error) {
	opts := &clientOptions{}
	for _, opt := range options {
		if err := opt(opts); err != nil {
			return nil, err
		}
	}
	if opts.domain == "" {
		return nil, fmt.Errorf("no domain specified")
	}

	dialCtx, cancelDial := context.WithTimeout(ctx, dialTimeout)
	defer cancelDial()
	authConn, err := grpc.DialContext(dialCtx, urlFromDomain("auth", opts.domain), opts.dialOptions...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	auth := pb_auth.NewAuthClient(authConn)

	togglesConn, err := grpc.DialContext(dialCtx, urlFromDomain("toggles", opts.domain), opts.dialOptions...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	toggles := pb_toggles.NewTogglesClient(togglesConn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cl := &Client{
		authClient:    auth,
		togglesClient: toggles,
		apiKey:        opts.apiKey,
	}

	return cl, nil
}

func urlFromDomain(service, domain string) string {
	return service + "." + domain
}
