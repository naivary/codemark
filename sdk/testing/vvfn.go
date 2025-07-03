package testing

import (
	"math"
	"reflect"

	sdkutil "github.com/naivary/codemark/sdk/utils"
)

func GetVVFn(rtype reflect.Type) ValidValueFunc {
	if sdkutil.IsValidSlice(rtype) {
		return isValidList
	}
	if sdkutil.IsInt(rtype) || sdkutil.IsUint(rtype) {
		return isValidInteger
	}
	if sdkutil.IsFloat(rtype) {
		return isValidFloat
	}
	if sdkutil.IsComplex(rtype) {
		return isValidComplex
	}
	if sdkutil.IsBool(rtype) {
		return isValidBool
	}
	if sdkutil.IsString(rtype) {
		return isValidString
	}
	return nil
}

func isValidInteger(got, want reflect.Value) bool {
	got = derefValue(got)
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

func isValidList(got, want reflect.Value) bool {
	elem := got.Type().Elem()
	vvfn := GetVVFn(elem)
	if vvfn == nil {
		return false
	}
	for i := range want.Len() {
		wantElem := want.Index(i)
		gotElem := got.Index(i)
		if !vvfn(gotElem, wantElem) {
			return false
		}
	}
	return true
}

func isValidFloat(got, want reflect.Value) bool {
	got = derefValue(got)
	w := want.Interface().(float64)
	return almostEqual(got.Float(), w)
}

func isValidBool(got, want reflect.Value) bool {
	got = derefValue(got)
	w := want.Interface().(bool)
	return got.Bool() == w
}

func isValidString(got, want reflect.Value) bool {
	got = derefValue(got)
	w := want.Interface().(string)
	return got.String() == w
}

func isValidComplex(got, want reflect.Value) bool {
	got = derefValue(got)
	w := want.Interface().(complex128)
	return got.Complex() == w
}

func almostEqual(a, b float64) bool {
	const float64EqualityThreshold = 1e-5
	return math.Abs(a-b) <= float64EqualityThreshold
}

func derefValue(v reflect.Value) reflect.Value {
	if sdkutil.IsPointer(v.Type()) {
		return v.Elem()
	}
	return v
}
