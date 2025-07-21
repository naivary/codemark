package converter

import (
	"reflect"
	"testing"

	coreapi "github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/converter/convertertest"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/registry/registrytest"
)

func customTests(tester convertertest.Tester) []convertertest.Case {
	b := marker.New("codemark:testing:byte", marker.STRING, reflect.ValueOf("b"))
	r := marker.New("codemark:testing:rune", marker.STRING, reflect.ValueOf("r"))
	return []convertertest.Case{
		tester.MustNewCaseWithMarker(&b, reflect.TypeFor[registrytest.String](), true, coreapi.TargetAny),
		tester.MustNewCaseWithMarker(&r, reflect.TypeFor[registrytest.String](), true, coreapi.TargetAny),
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
	tester, err := newConvTester(conv, registrytest.StringTypes())
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range customTests(tester) {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}
