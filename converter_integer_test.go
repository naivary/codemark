package codemark

import (
	"math"
	"reflect"
	"slices"
	"testing"

	sdktesting "github.com/naivary/codemark/sdk/testing"
	"github.com/naivary/codemark/sdk/utils"
)

func isValidInteger(got, want reflect.Value) bool {
	got = utils.DeRefValue(got)
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

func newConvTesterForInteger() (sdktesting.ConverterTester, error) {
	tester, err := sdktesting.NewConverterTester(&intConverter{})
	if err != nil {
		return nil, err
	}
	for _, typ := range slices.Concat(sdktesting.IntTypes(), sdktesting.UintTypes()) {
		rtype := reflect.TypeOf(typ)
		if err := tester.AddType(rtype); err != nil {
			return nil, err
		}
		if err := tester.AddVVFunc(rtype, isValidInteger); err != nil {
			return nil, err
		}
	}
	return tester, nil
}

func TestIntConverter(t *testing.T) {
	reg, err := sdktesting.NewDefsSet(NewInMemoryRegistry(), &DefinitionMarker{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	mngr, err := NewConvMngr(reg)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tester, err := newConvTesterForInteger()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tests, err := tester.ValidTests()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range tests {
		tester.Run(t, tc, mngr)
	}
}
