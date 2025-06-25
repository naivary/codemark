package codemark

import (
	"math"
	"reflect"
	"slices"

	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

func newConvTester(conv sdk.Converter) (sdktesting.ConverterTester, error) {
	tester, err := sdktesting.NewConverterTester(conv, nil)
	if err != nil {
		return nil, err
	}
	types := getTypes(conv)
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		vvfn := getVVFn(rtype)
		if err := tester.AddType(rtype); err != nil {
			return nil, err
		}
		if err := tester.AddVVFunc(rtype, vvfn); err != nil {
			return nil, err
		}
	}
	return tester, nil
}

func getTypes(conv sdk.Converter) []any {
	if _, isList := conv.(*listConverter); isList {
		return sdktesting.ListTypes()
	}
	if _, isInt := conv.(*intConverter); isInt {
		return slices.Concat(sdktesting.IntTypes(), sdktesting.UintTypes())
	}
	if _, isFloat := conv.(*floatConverter); isFloat {
		return sdktesting.FloatTypes()
	}
	if _, isComplex := conv.(*complexConverter); isComplex {
		return sdktesting.ComplexTypes()
	}
	if _, isBool := conv.(*boolConverter); isBool {
		return sdktesting.BoolTypes()
	}
	if _, isString := conv.(*stringConverter); isString {
		return sdktesting.StringTypes()
	}
	return nil

}

func getVVFn(rtype reflect.Type) sdktesting.ValidValueFunc {
	typeID := sdkutil.TypeIDOf(rtype)
	if sdkutil.MatchTypeID(typeID, `slice\..+`) {
		return isValidList
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?int\d{0,2}`) {
		return isValidInteger
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?float\d{2}`) {
		return isValidFloat
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?complex\d{2,3}`) {
		return isValidComplex
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?bool`) {
		return isValidBool
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?string`) {
		return isValidString
	}
	return nil
}

func isValidInteger(got, want reflect.Value) bool {
	got = sdkutil.DeRefValue(got)
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
	vvfn := getVVFn(elem)
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
	got = sdkutil.DeRefValue(got)
	w := want.Interface().(float64)
	return sdktesting.AlmostEqual(got.Float(), w)
}

func isValidBool(got, want reflect.Value) bool {
	got = sdkutil.DeRefValue(got)
	w := want.Interface().(bool)
	return got.Bool() == w
}

func isValidString(got, want reflect.Value) bool {
	got = sdkutil.DeRefValue(got)
	w := want.Interface().(string)
	return got.String() == w
}

func isValidComplex(got, want reflect.Value) bool {
	got = sdkutil.DeRefValue(got)
	w := want.Interface().(complex128)
	return got.Complex() == w
}
