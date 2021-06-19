package factory

import (
	"github.com/imdario/mergo"
	"github.com/tjpnz/structbot/pkg/config"
	"testing"
)

// Factory defines methods to create and patch arbitrary structs.
type Factory interface {
	Create() (interface{}, error)
	Patch(v interface{}) (interface{}, error)
	MustCreate(t testing.TB) interface{}
	MustPatch(t testing.TB, v interface{}) interface{}
}

type factoryImpl struct {
	config *config.Config
	fn     func() (interface{}, error)
}

// New returns a new Factory.
func New(config *config.Config, fn func() (interface{}, error)) Factory {
	return &factoryImpl{
		config: config,
		fn: fn,
	}
}

// Create invokes the registered factory function. It returns a new struct or an
// error if it couldn't be created.
func (f *factoryImpl) Create() (interface{}, error) {
	return f.fn()
}

// Patch invokes the registered factory function and then updates it with the
// fields in v. If a field is missing in v the field on the base struct is
// not overwritten. If the struct couldn't be patched an error is returned.
func (f *factoryImpl) Patch(v interface{}) (interface{}, error) {
	base, err := f.fn()
	if err != nil {
		return nil, err
	}

	if err := mergo.Merge(base, v, mergo.WithOverride, mergo.WithTransformers(f.config.Transformers)); err != nil {
		return nil, err
	}
	return base, nil
}

// MustCreate is a variant of Create but calls t.Fatalf on error.
func (f *factoryImpl) MustCreate(t testing.TB) interface{} {
	out, err := f.fn()
	if err != nil {
		t.Fatalf("failed to create struct: %v", err)
	}
	return out
}

// MustPatch is a variant of Patch but calls t.Fatalf on error.
func (f *factoryImpl) MustPatch(t testing.TB, v interface{}) interface{} {
	out, err := f.Patch(v)
	if err != nil {
		t.Fatalf("failed to patch struct: %v", err)
	}
	return out
}
