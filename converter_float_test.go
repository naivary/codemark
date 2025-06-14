package codemark

import (
	"reflect"
	"testing"

	sdktesting "github.com/naivary/codemark/sdk/testing"
	"github.com/naivary/codemark/sdk/utils"
)

func isValidFloat(got, want reflect.Value) bool {
	got = utils.DeRefValue(got)
	w := want.Interface().(float64)
	return sdktesting.AlmostEqual(got.Float(), w)
}

func newConvTesterForFloat() (sdktesting.ConverterTester, error) {
	tester, err := sdktesting.NewConverterTester(&floatConverter{})
	if err != nil {
		return nil, err
	}
	for _, typ := range sdktesting.FloatTypes() {
		rtype := reflect.TypeOf(typ)
		if err := tester.AddType(rtype); err != nil {
			return nil, err
		}
		if err := tester.AddVVFunc(rtype, isValidFloat); err != nil {
			return nil, err
		}
	}
	return tester, nil
}

func TestFloatConverter(t *testing.T) {
	reg, err := sdktesting.NewDefsSet(NewInMemoryRegistry(), &DefinitionMarker{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	mngr, err := NewConvMngr(reg)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tester, err := newConvTesterForFloat()
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
