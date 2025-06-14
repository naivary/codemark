package codemark

import (
	"reflect"
	"testing"

	sdktesting "github.com/naivary/codemark/sdk/testing"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

func getElemVVFn(typeID string) sdktesting.ValidValueFunc {
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

func isValidList() sdktesting.ValidValueFunc {
	vvfn := func(got, want reflect.Value) bool {
		elem := got.Type().Elem()
		typeID := sdkutil.TypeIDOf(elem)
		vvfn := getElemVVFn(typeID)
		if vvfn == nil {
			return false
		}
		for i := 0; i < want.Len(); i++ {
			wantElem := want.Index(i)
			gotElem := got.Index(i)
			if !vvfn(gotElem, wantElem) {
				return false
			}
		}
		return true
	}
	return vvfn
}

func newConvTesterForList() (sdktesting.ConverterTester, error) {
	tester, err := sdktesting.NewConverterTester(&listConverter{})
	if err != nil {
		return nil, err
	}
	for _, typ := range sdktesting.ListTypes() {
		rtype := reflect.TypeOf(typ)
		if err := tester.AddType(rtype); err != nil {
			return nil, err
		}
		if err := tester.AddVVFunc(rtype, isValidList()); err != nil {
			return nil, err
		}
	}
	return tester, nil
}

func TestListConverter(t *testing.T) {
	reg, err := sdktesting.NewDefsSet(NewInMemoryRegistry(), &DefinitionMarker{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	mngr, err := NewConvMngr(reg)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tester, err := newConvTesterForList()
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
