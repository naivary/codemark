package converter

import (
	"reflect"
	"slices"

	coreapi "github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func newConvTester(conv sdk.Converter, customTypes []any) (sdktesting.ConverterTester, error) {
	tester, err := sdktesting.NewConvTester(conv)
	if err != nil {
		return nil, err
	}
	for _, typ := range customTypes {
		to := reflect.TypeOf(typ)
		vvfn := sdktesting.GetVVFn(to)
		if err := tester.AddVVFunc(to, vvfn); err != nil {
			return nil, err
		}
	}
	return tester, nil
}

func customTypesFor(conv sdk.Converter) []any {
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

func validTestsFor(conv sdk.Converter, tester sdktesting.ConverterTester) ([]sdktesting.ConverterTestCase, error) {
	types := customTypesFor(conv)
	tests := make([]sdktesting.ConverterTestCase, 0, len(types))
	for _, to := range types {
		rtype := reflect.TypeOf(to)
		tc, err := tester.NewTest(rtype, true, coreapi.TargetAny)
		if err != nil {
			return nil, err
		}
		tests = append(tests, tc)
	}
	return tests, nil
}
