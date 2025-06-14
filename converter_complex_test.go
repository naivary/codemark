package codemark

import (
	"reflect"
	"testing"

	sdktesting "github.com/naivary/codemark/sdk/testing"
	"github.com/naivary/codemark/sdk/utils"
)

func isValidComplex(got, want reflect.Value) bool {
	got = utils.DeRefValue(got)
	w := want.Interface().(complex128)
	return got.Complex() == w
}

func newConvTesterForComplex() (sdktesting.ConverterTester, error) {
	tester, err := sdktesting.NewConverterTester(&complexConverter{})
	if err != nil {
		return nil, err
	}
	for _, typ := range sdktesting.ComplexTypes() {
		rtype := reflect.TypeOf(typ)
		if err := tester.AddType(rtype); err != nil {
			return nil, err
		}
		if err := tester.AddVVFunc(rtype, isValidComplex); err != nil {
			return nil, err
		}
	}
	return tester, nil
}

func TestComplexConverter(t *testing.T) {
	reg, err := sdktesting.NewDefsSet(NewInMemoryRegistry(), &DefinitionMarker{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	mngr, err := NewConvMngr(reg)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tester, err := newConvTesterForComplex()
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
