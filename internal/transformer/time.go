package transformer

import (
	"reflect"
	"time"
)

// TimeTransformer implements a transformer for additional handling of
// time.Time structs.
type TimeTransformer struct {}

// Transformer handles zeroed time.Time structs.
func (t TimeTransformer) Transformer(typ reflect.Type) func (dst, src reflect.Value) error {
	if typ == reflect.TypeOf(time.Time{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				if !isZero(src) {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	return nil
}

func isZero(v reflect.Value) bool {
	m := v.MethodByName("IsZero")
	res := m.Call([]reflect.Value{})
	return res[0].Bool()
}
