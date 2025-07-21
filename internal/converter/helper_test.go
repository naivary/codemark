package converter

import (
	"reflect"
	"slices"

	coreapi "github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/converter/convertertest"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/registry/registrytest"
)

func newConvTester(conv converter.Converter, customTypes []any) (convertertest.Tester, error) {
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		return nil, err
	}
	for _, typ := range customTypes {
		to := reflect.TypeOf(typ)
		equal := marker.GetEqualFunc(to)
		if err := tester.AddEqualFunc(to, equal); err != nil {
			return nil, err
		}
	}
	return tester, nil
}

func customTypesFor(conv converter.Converter) []any {
	if _, isList := conv.(*listConverter); isList {
		return registrytest.ListTypes()
	}
	if _, isInt := conv.(*intConverter); isInt {
		return slices.Concat(registrytest.IntTypes(), registrytest.UintTypes())
	}
	if _, isFloat := conv.(*floatConverter); isFloat {
		return registrytest.FloatTypes()
	}
	if _, isComplex := conv.(*complexConverter); isComplex {
		return registrytest.ComplexTypes()
	}
	if _, isBool := conv.(*boolConverter); isBool {
		return registrytest.BoolTypes()
	}
	if _, isString := conv.(*stringConverter); isString {
		return registrytest.StringTypes()
	}
	return nil
}

func validTestsFor(conv converter.Converter, tester convertertest.Tester) ([]convertertest.Case, error) {
	types := customTypesFor(conv)
	tests := make([]convertertest.Case, 0, len(types))
	for _, to := range types {
		rtype := reflect.TypeOf(to)
		tc, err := tester.NewCase(rtype, true, coreapi.TargetAny)
		if err != nil {
			return nil, err
		}
		tests = append(tests, tc)
	}
	return tests, nil
}
