package featureguards

import "time"

// WithAttributes specifies which attributes are passed to FeatureGuards for rule evaluation.
// For example, user_id, session_id or other attributes. Note: It's case sensitive.
func WithAttributes(a Attributes) FeatureToggleOptions {
	return func(o *ftOptions) error {
		o.attrs = a
		return nil
	}
}

// Attributes is a dictionary of keys that will be used for evaluation to values.
type Attributes map[string]interface{}

// Int64 adds a new int64 attribute.
func (a Attributes) Int64(name string, n int64) Attributes {
	a[name] = int64(n)
	return a
}

// Int adds a new int attribute.
func (a Attributes) Int(name string, n int) Attributes {
	a[name] = int64(n)
	return a
}

// Float adds a new float32 attribute.
func (a Attributes) Float(name string, n float32) Attributes {
	a[name] = n
	return a
}

// String adds a new string attribute.
func (a Attributes) String(name string, v string) Attributes {
	a[name] = v
	return a
}

// Bool adds a new boolean attribute.
func (a Attributes) Bool(name string, v bool) Attributes {
	a[name] = v
	return a
}

// Time adds a new time.Time attribute.
func (a Attributes) Time(name string, t time.Time) Attributes {
	a[name] = t
	return a
}
