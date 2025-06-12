package testing

import (
	"reflect"
)

// IsValidValue checks if the wanted value is equest to the value of got after
// converting it to the type `T`. If your produced got value is expected to be a
// slice use `slices.Equal` for comparison.
func IsValidValue[T comparable](got reflect.Value, wanted T) bool {
	v := got.Interface().(T)
	return v == wanted
}
