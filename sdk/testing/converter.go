package testing

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
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


func IsValidValuePtrFloat[T ~*float32 | ~*float64](got, wanted reflect.Value) bool {
	// got is a reflect.Value containing a *float32 or *float64
	ptr := got.Interface().(T) // T is *float32 or *float64
	value := reflect.ValueOf(ptr).Elem().Interface()
	var f64 float64
	switch v := value.(type) {
	case float32:
		f64 = float64(v)
	case float64:
		f64 = v
	default:
		return false
	}
	w := wanted.Interface().(float64)
	return AlmostEqual(f64, w)
}
