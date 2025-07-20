package converter

import (
	"reflect"
	"testing"

	coreapi "github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/parser/marker"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func customTests(tester sdktesting.ConverterTester) []sdktesting.ConverterTestCase {
	b := marker.New("codemark:testing:byte", marker.STRING, reflect.ValueOf("b"))
	r := marker.New("codemark:testing:rune", marker.STRING, reflect.ValueOf("r"))
	return []sdktesting.ConverterTestCase{
		tester.MustNewTestWithMarker(&b, reflect.TypeFor[sdktesting.String](), true, coreapi.TargetAny),
		tester.MustNewTestWithMarker(&r, reflect.TypeFor[sdktesting.String](), true, coreapi.TargetAny),
	}
}

func TestIntConverter(t *testing.T) {
	conv := Integer()
	tester, err := newConvTester(conv, customTypesFor(conv))
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tests, err := validTestsFor(conv, tester)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}

func TestIntConverter_Byte_and_Rune(t *testing.T) {
	conv := Integer()
	tester, err := newConvTester(conv, sdktesting.StringTypes())
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range customTests(tester) {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}
