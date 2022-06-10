package featureguards

import (
	"github.com/featureguards/featureguards-go/v2/internal/logger"

	"google.golang.org/grpc"
)

type LogLevel = logger.LogLevel
type Options = func(o *toggleOptions) error

type toggleOptions struct {
	domain        string
	logLevel      logger.LogLevel
	dialOptions   []grpc.DialOption
	testCerts     bool
	apiKey        string
	withoutListen bool
	defaults      map[string]bool
}

func WithDialOptions(options ...grpc.DialOption) Options {
	return func(o *toggleOptions) error {
		o.dialOptions = options
		return nil
	}
}

func WithApiKey(key string) Options {
	return func(o *toggleOptions) error {
		o.apiKey = key
		return nil
	}
}

func WithDefaults(v map[string]bool) Options {
	return func(o *toggleOptions) error {
		o.defaults = v
		return nil
	}
}

// For internal testing purposes mostly.
func withDomain(domain string) Options {
	return func(o *toggleOptions) error {
		o.domain = domain
		return nil
	}
}

func withoutListen() Options {
	return func(o *toggleOptions) error {
		o.withoutListen = true
		return nil
	}
}

func withTestCerts() Options {
	return func(o *toggleOptions) error {
		o.testCerts = true
		return nil
	}
}
