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
				isZero := dst.MethodByName("IsZero")
				result := isZero.Call([]reflect.Value{})
				if result[0].Bool() {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	return nil
}
