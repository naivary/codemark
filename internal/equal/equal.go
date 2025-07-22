package equal

import (
	"math"
	"reflect"

	"github.com/naivary/codemark/typeutil"
)

func IsEqual(got, want reflect.Value) bool {
	equal := GetFunc(got.Type())
	if equal == nil {
		return false
	}
	return equal(got, want)
}

// GetFunc returns a function which can be used to compare to values. One
// being the marker value and the other a value produced by some processing.
func GetFunc(rtype reflect.Type) func(got, want reflect.Value) bool {
	if typeutil.IsValidSlice(rtype) {
		return list
	}
	if typeutil.IsInt(rtype) || typeutil.IsUint(rtype) {
		return integer
	}
	if typeutil.IsFloat(rtype) {
		return float
	}
	if typeutil.IsComplex(rtype) {
		return complexx
	}
	if typeutil.IsBool(rtype) {
		return boolean
	}
	if typeutil.IsString(rtype) {
		return stringg
	}
	if typeutil.IsAny(rtype) {
		return anything
	}
	return nil
}

func list(got, want reflect.Value) bool {
	for i := range want.Len() {
		wantElem := want.Index(i)
		gotElem := got.Index(i)
		equal := GetFunc(gotElem.Type())
		if equal == nil {
			return false
		}
		if !equal(gotElem, wantElem) {
			return false
		}
	}
	return true
}

func anything(got, want reflect.Value) bool {
	got = typeutil.DerefValue(got)
	want = typeutil.DerefValue(want)
	return reflect.DeepEqual(got.Interface(), want.Interface())
}

func integer(got, want reflect.Value) bool {
	got = typeutil.DerefValue(got)
	var i64 int64
	switch got.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 = got.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := got.Uint()
		if u64 > math.MaxInt64 {
			return false // uint64 value too large for int64
		}
		i64 = int64(u64)
	default:
		return false
	}
	w := want.Interface().(int64)
	return i64 == w
}

func float(got, want reflect.Value) bool {
	got = typeutil.DerefValue(got)
	w := want.Interface().(float64)
	return almostEqual(got.Float(), w)
}

func boolean(got, want reflect.Value) bool {
	got = typeutil.DerefValue(got)
	w := want.Interface().(bool)
	return got.Bool() == w
}

func stringg(got, want reflect.Value) bool {
	got = typeutil.DerefValue(got)
	w := want.Interface().(string)
	return got.String() == w
}

func complexx(got, want reflect.Value) bool {
	got = typeutil.DerefValue(got)
	w := want.Interface().(complex128)
	return got.Complex() == w
}

func almostEqual(a, b float64) bool {
	const float64EqualityThreshold = 1e-5
	return math.Abs(a-b) <= float64EqualityThreshold
}
