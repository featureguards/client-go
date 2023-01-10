package dynamic_settings

import (
	"errors"
	"sync"

	pb_ds "github.com/featureguards/featureguards-go/v2/proto/dynamic_setting"

	log "github.com/sirupsen/logrus"
)

var (
	ErrDuplicate  error = errors.New("duplicate dynamic setting")
	ErrNoVariable error = errors.New("no variable passed")
)

type DynamicSettings struct {
	mu sync.RWMutex

	settings map[string]*setting
}

type BoolOption func(o *options)
type Int64Option func(o *options)
type Float32Option func(o *options)
type StringOption func(o *options)

type options struct {
	boolV    bool
	int64V   int64
	float32V float32
	stringV  string
}

type setting struct {
	serverCopy *pb_ds.DynamicSetting
	passedType pb_ds.DynamicSetting_Type
	passedVar  interface{}
}

// New returns a new dynamic settings wrapper to be used to define typed dynamic settings that
// are updated automatically via feature guards.
func New() *DynamicSettings {
	return &DynamicSettings{
		settings: make(map[string]*setting),
	}
}

// WithDefaultBool defines a default bool value to be used.
func WithDefaultBool(v bool) BoolOption {
	return func(o *options) {
		o.boolV = v
	}
}

// WithDefaultInt64 defines a default int64 value to be used.
func WithDefaultInt64(v int64) Int64Option {
	return func(o *options) {
		o.int64V = v
	}
}

// WithDefaultFloat32 defines a default float32 value to be used.
func WithDefaultFloat32(v float32) Float32Option {
	return func(o *options) {
		o.float32V = v
	}
}

// WithDefaultString defines a default string value to be used.
func WithDefaultString(v string) StringOption {
	return func(o *options) {
		o.stringV = v
	}
}

// Bool defines a bool dynamic setting. A default value can be passed optionally.
func (ds *DynamicSettings) Bool(name string, v *bool, opts ...BoolOption) error {
	if v == nil {
		return ErrNoVariable
	}
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()
	_, ok := ds.settings[name]
	if ok {
		return ErrDuplicate
	}

	*v = o.boolV
	ds.settings[name] = &setting{
		passedVar:  v,
		passedType: pb_ds.DynamicSetting_BOOL,
	}
	return nil
}

// Int64 defines a int64 dynamic setting. A default value can be passed optionally.
func (ds *DynamicSettings) Int64(name string, v *int64, opts ...Int64Option) error {
	if v == nil {
		return ErrNoVariable
	}
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()
	_, ok := ds.settings[name]
	if ok {
		return ErrDuplicate
	}

	*v = o.int64V
	ds.settings[name] = &setting{
		passedVar:  v,
		passedType: pb_ds.DynamicSetting_INTEGER,
	}
	return nil
}

// Float32 defines a float32 dynamic setting. A default value can be passed optionally.
func (ds *DynamicSettings) Float32(name string, v *float32, opts ...Float32Option) error {
	if v == nil {
		return ErrNoVariable
	}
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()
	_, ok := ds.settings[name]
	if ok {
		return ErrDuplicate
	}

	*v = o.float32V
	ds.settings[name] = &setting{
		passedVar:  v,
		passedType: pb_ds.DynamicSetting_FLOAT,
	}
	return nil
}

// String defines a string dynamic setting. A default value can be passed optionally.
func (ds *DynamicSettings) String(name string, v *string, opts ...StringOption) error {
	if v == nil {
		return ErrNoVariable
	}
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()
	_, ok := ds.settings[name]
	if ok {
		return ErrDuplicate
	}

	*v = o.stringV
	ds.settings[name] = &setting{
		passedVar:  v,
		passedType: pb_ds.DynamicSetting_STRING,
	}
	return nil
}

// Process is used internally to update the settings based on a newer set.
func (ds *DynamicSettings) Process(settings []*pb_ds.DynamicSetting) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for _, setting := range settings {
		if setting.DeletedAt.IsValid() {
			delete(ds.settings, setting.Name)
		} else {
			found := ds.settings[setting.Name]
			if found == nil {
				// Skip unknown variable
				continue
			}
			if found.passedType != setting.SettingType {
				log.Errorf("Dynamic setting type: %v != %v for %s. Skipping.\n", found.passedType, setting.SettingType, setting.Name)
			}
			found.serverCopy = setting

			// update stored variable. We've already validated the type of the variables.
			switch v := found.passedVar.(type) {
			case *bool:
				*v = setting.GetBoolValue().Value
			case *int64:
				*v = setting.GetIntegerValue().Value
			case *float32:
				*v = setting.GetFloatValue().Value
			case *string:
				*v = setting.GetStringValue().Value
			default:
				log.Warningf("Skipping unknown type for %s.\n", setting.Name)
			}
		}
	}
	return nil
}
