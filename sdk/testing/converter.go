package testing

import (
	"fmt"
	"math"
	"reflect"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
	"golang.org/x/exp/constraints"
)

type ConverterTestCase struct {
	Name         string
	Marker       parser.Marker
	Target       sdk.Target
	ToType       reflect.Type
	IsValidCase  bool
	IsValidValue func(got reflect.Value, wanted reflect.Value) bool
}

func NewConvTestCases(
	conv sdk.Converter,
	validValueFuncs map[string]func(got reflect.Value, wanted reflect.Value) bool,
) ([]ConverterTestCase, error) {
	tests := make([]ConverterTestCase, 0, len(conv.SupportedTypes()))
	for _, rtype := range conv.SupportedTypes() {
		typeID := sdkutil.TypeID(rtype)
		marker := RandMarkerFromRefType(rtype)
		if marker == nil {
			return nil, fmt.Errorf("no valid marker found: %v\n", rtype)
		}
		tc := ConverterTestCase{
			Name:         "random-name",
			Marker:       marker,
			Target:       sdk.TargetAny,
			ToType:       rtype,
			IsValidCase:  true,
			IsValidValue: validValueFuncs[typeID],
		}
		tests = append(tests, tc)
	}
	return tests, nil
}

func IsValidInteger[T constraints.Integer](got, wanted reflect.Value) bool {
	g := got.Interface().(T)
	w := wanted.Interface().(int64)
	return int64(g) == w
}

func IsValidValuePtr[T ~*int | ~*uint](got, wanted reflect.Value) bool {
	ptr := got.Interface().(T)
	value := reflect.ValueOf(ptr).Elem()
	var i64 int64
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 = value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := value.Uint()
		if u64 > math.MaxInt64 {
			return false // uint64 value too large for int64
		}
	default:
		return false
	}
	w := wanted.Interface().(int64)
	return i64 == w
}

func IsValidValuePtrFloat[T ~*float32 | ~*float64](got, wanted reflect.Value) bool {
	ptr := got.Interface().(T)
	value := reflect.ValueOf(ptr).Elem().Float()
	w := wanted.Interface().(float64)
	return AlmostEqual(value, w)
}
