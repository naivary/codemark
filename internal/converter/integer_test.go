package converter

import (
	"reflect"
	"testing"
	"time"

	"github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/converter/convertertest"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/marker/markertest"
	"github.com/naivary/codemark/typeutil"
)

type (
	Duration    time.Duration
	PtrDuration *time.Duration
)

func equalDuration(got, want reflect.Value) bool {
	t := reflect.TypeFor[time.Duration]()
	gotDuration := typeutil.DerefValue(got).Convert(t).Interface().(time.Duration)
	wantDuration, err := time.ParseDuration(want.String())
	if err != nil {
		return false
	}
	return gotDuration == wantDuration
}

// TODO: byte and rune test cases missing
// TODO: target cann raus beim convertertest.Tester
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

func TestIntConverter_Duration(t *testing.T) {
	conv := Integer()
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	equalFuncs := map[reflect.Type]func(got, want reflect.Value) bool{
		reflect.TypeFor[Duration]():    equalDuration,
		reflect.TypeFor[PtrDuration](): equalDuration,
	}
	for typ, equalFunc := range equalFuncs {
		if err := tester.AddEqualFunc(typ, equalFunc); err != nil {
			t.Errorf("err occured: %s\n", err)
		}
	}
	duration := markertest.NewMarker("duration", marker.STRING, reflect.ValueOf("10h"))
	ptrDuration := markertest.NewMarker("ptr.duration", marker.STRING, reflect.ValueOf("10h"))
	cases := []convertertest.Case{
		tester.MustNewCaseWithMarker(&duration, reflect.TypeOf(Duration(0)), true, core.TargetAny),
		tester.MustNewCaseWithMarker(&ptrDuration, reflect.TypeOf(PtrDuration(nil)), true, core.TargetAny),
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}
