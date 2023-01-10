package featureguards

import (
	"github.com/featureguards/featureguards-go/v2/dynamic_settings"
	"github.com/featureguards/featureguards-go/v2/internal/logger"

	"google.golang.org/grpc"
)

// LogLevel controls the verbosity logging level.
type LogLevel = logger.LogLevel

// Options specifies the options passed to the FeatueGuards client.
type Options = func(o *toggleOptions) error

type toggleOptions struct {
	domain          string
	logLevel        logger.LogLevel
	dialOptions     []grpc.DialOption
	testCerts       bool
	apiKey          string
	withoutListen   bool
	defaults        map[string]bool
	dynamicSettings *dynamic_settings.DynamicSettings
}

// WithDialOptions adds gRPC dial options. Can be used to enforce a timeout on the initial dial.
func WithDialOptions(options ...grpc.DialOption) Options {
	return func(o *toggleOptions) error {
		o.dialOptions = options
		return nil
	}
}

// WithApiKey is required and specifies the API key specific to the FeatureGuards project and environment.
func WithApiKey(key string) Options {
	return func(o *toggleOptions) error {
		o.apiKey = key
		return nil
	}
}

// WithDefaults adds default values for feature toggle names. This is useful to ensure that in
// cases where FeatureGuards is down or cannot be reached, you can specify different values to be
// returned. By default, every feature toggle is off unless a different value is specified here.
func WithDefaults(v map[string]bool) Options {
	return func(o *toggleOptions) error {
		o.defaults = v
		return nil
	}
}

// WithDynamicSettings adds the dynamic settings to be updated via FeatureGuards.
func WithDynamicSettings(v *dynamic_settings.DynamicSettings) Options {
	return func(o *toggleOptions) error {
		o.dynamicSettings = v
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
