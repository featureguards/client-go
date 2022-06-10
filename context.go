package featureguards

import "time"

func WithAttributes(a Attributes) FeatureToggleOptions {
	return func(o *ftOptions) error {
		o.attrs = a
		return nil
	}
}

type Attributes map[string]interface{}

func (a Attributes) Int64(name string, n int64) Attributes {
	a[name] = float32(n)
	return a
}

func (a Attributes) Int(name string, n int) Attributes {
	a[name] = float32(n)
	return a
}

func (a Attributes) Float(name string, n float32) Attributes {
	a[name] = n
	return a
}

func (a Attributes) String(name string, v string) Attributes {
	a[name] = v
	return a
}

func (a Attributes) Bool(name string, v bool) Attributes {
	a[name] = v
	return a
}

func (a Attributes) Time(name string, t time.Time) Attributes {
	a[name] = t
	return a
}
