package converter

import (
	"reflect"
	"slices"
	"time"

	coreapi "github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/converter/convertertest"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/registry/registrytest"
	"github.com/naivary/codemark/typeutil"
)

func newConvTester(conv converter.Converter, customTypes []any) (convertertest.Tester, error) {
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		return nil, err
	}
	for _, typ := range customTypes {
		to := reflect.TypeOf(typ)
		equal := getEqualFunc(to)
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
		return slices.Concat(registrytest.StringTypes())
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

func getEqualFunc(to reflect.Type) func(got, want reflect.Value) bool {
	derefed := typeutil.Deref(to)
	if derefed == reflect.TypeOf(Time(time.Time{})) || derefed == reflect.TypeOf(time.Time{}) {
		return equalTime
	}
	return marker.GetEqualFunc(to)
}

func equalTime(got, want reflect.Value) bool {
	t := reflect.TypeOf(time.Time{})
	gotTime := typeutil.DerefValue(got).Convert(t).Interface().(time.Time)
	wantTime, err := time.Parse(time.RFC3339, want.String())
	if err != nil {
		return false
	}
	return gotTime.Equal(wantTime)
}
