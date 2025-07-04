package converter

import (
	"reflect"
	"slices"

	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/registry"
	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

var reg = newRegistry()

var mngr = newManager()

func newManager() *Manager {
	mngr, err := NewManager(reg)
	if err != nil {
		panic(err)
	}
	return mngr
}

func newRegistry() registry.Registry {
	defs := sdktesting.NewDefSet()
	reg, err := sdktesting.NewRegistry(defs)
	if err != nil {
		panic(err)
	}
	return reg
}

func newConvTester(conv sdk.Converter, customeTypes []any) (sdktesting.ConverterTester, error) {
	tester, err := sdktesting.NewConvTester(conv)
	if err != nil {
		return nil, err
	}
	for _, typ := range customeTypes {
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
	// if _, isAny := conv.(*anyConverter); isAny {
	// 	return sdktesting.AnyTypes()
	// }
	return nil
}

func validTestsFor(conv sdk.Converter, tester sdktesting.ConverterTester) ([]sdktesting.ConverterTestCase, error) {
	types := customTypesFor(conv)
	tests := make([]sdktesting.ConverterTestCase, 0, len(types))
	for _, to := range types {
		rtype := reflect.TypeOf(to)
		tc, err := tester.NewTest(rtype, true, target.ANY)
		if err != nil {
			return nil, err
		}
		tests = append(tests, tc)
	}
	return tests, nil
}
