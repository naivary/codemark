package testing

import (
	"fmt"
	"math"
	"reflect"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	"github.com/naivary/codemark/sdk/utils"
	sdkutil "github.com/naivary/codemark/sdk/utils"
	"golang.org/x/exp/constraints"
)

type ValidValueFunc func(got, want reflect.Type) bool

type ConverterTester interface {
	GetValidValueFunc(typeID string) ValidValueFunc
}

var validValueFuncs = map[string]func(got, want reflect.Value) bool{}

type ConverterTestCase struct {
	Name         string
	Marker       parser.Marker
	Target       sdk.Target
	ToType       reflect.Type
	IsValidCase  bool
	IsValidValue func(got reflect.Value, wanted reflect.Value) bool
}

func NewConvTestCases(conv sdk.Converter) ([]ConverterTestCase, error) {
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

func getValidValueFunc(typeID string) func(got, want reflect.Value) bool {
	return nil
}

func isValidInteger[T constraints.Integer | ~*int | ~*int8 | ~*int16 | ~*int32 | ~*int64 | ~*uint | ~*uint8 | ~*uint16 | ~*uint32 | ~*uint64](got, want reflect.Value) bool {
	if utils.IsPointer(got.Type()) {
		got = got.Elem()
	}
	var i64 int64
	switch got.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 = got.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := got.Uint()
		if u64 > math.MaxInt64 {
			return false // uint64 value too large for int64
		}
	default:
		return false
	}
	w := want.Interface().(int64)
	return i64 == w
}
