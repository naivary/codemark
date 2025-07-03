package converter

import (
	"reflect"
	"slices"
	"testing"

	"github.com/naivary/codemark/definition/target"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func customTests(tester sdktesting.ConverterTester) []sdktesting.ConverterTestCase {
	return []sdktesting.ConverterTestCase{
		tester.MustNewTest(reflect.TypeFor[byte](), reflect.TypeFor[sdktesting.String](), true, target.ANY),
		tester.MustNewTest(reflect.TypeFor[*byte](), reflect.TypeFor[sdktesting.String](), true, target.ANY),
		tester.MustNewTest(reflect.TypeFor[rune](), reflect.TypeFor[sdktesting.String](), true, target.ANY),
		tester.MustNewTest(reflect.TypeFor[*rune](), reflect.TypeFor[sdktesting.String](), true, target.ANY),
	}
}

// ADD test for rune and byte e.g. from "s" to byte
func TestIntConverter(t *testing.T) {
	tester, err := newConvTester(&intConverter{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tests, err := tester.ValidTests()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range slices.Concat(tests, customTests(tester)) {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}

func TestIntConverter_String(t *testing.T) {
	tester, err := sdktesting.NewConvTester(&intConverter{}, mngr, map[reflect.Type]reflect.Type{
		reflect.TypeFor[byte]():  reflect.TypeFor[sdktesting.String](),
		reflect.TypeFor[*byte](): reflect.TypeFor[sdktesting.PtrString](),
		reflect.TypeFor[rune]():  reflect.TypeFor[sdktesting.String](),
		reflect.TypeFor[*byte](): reflect.TypeFor[sdktesting.PtrString](),
	})
	if err != nil {
		t.Fatalf("err occured: %s\n", err)
	}
	for _, typ := range sdktesting.StringTypes() {
		to := reflect.TypeOf(typ)
		vvfn := sdktesting.GetVVFn(to)
		if err := tester.AddVVFunc(to, vvfn); err != nil {
			t.Fatalf("err occured: %s\n", err)
		}
	}
	tests, err := tester.ValidTests()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range slices.Concat(tests, customTests(tester)) {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}

}
