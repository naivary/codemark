package converter

import (
	"reflect"
	"testing"
	"time"

	"github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/converter/convertertest"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/marker/markertest"
	"github.com/naivary/codemark/registry/registrytest"
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

func equalByte(got, want reflect.Value) bool {
	got = typeutil.DerefValue(got)
	return want.String() == string(byte(got.Uint()))
}

func equalRune(got, want reflect.Value) bool {
	got = typeutil.DerefValue(got)
	return want.String() == string(rune(got.Int()))
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
		tester.MustNewCaseWithMarker(
			&ptrDuration,
			reflect.TypeOf(PtrDuration(nil)),
			true,
			core.TargetAny,
		),
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}

func TestIntConverter_Byte(t *testing.T) {
	conv := Integer()
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	equalFuncs := map[reflect.Type]func(got, want reflect.Value) bool{
		reflect.TypeFor[registrytest.Byte]():    equalByte,
		reflect.TypeFor[registrytest.PtrByte](): equalByte,
	}
	for typ, equalFunc := range equalFuncs {
		if err := tester.AddEqualFunc(typ, equalFunc); err != nil {
			t.Errorf("err occured: %s\n", err)
		}
	}
	b := markertest.NewMarker("byte", marker.STRING, reflect.ValueOf("b"))
	ptrByte := markertest.NewMarker("ptr.byte", marker.STRING, reflect.ValueOf("b"))
	cases := []convertertest.Case{
		tester.MustNewCaseWithMarker(&b, reflect.TypeOf(registrytest.Byte(0)), true, core.TargetAny),
		tester.MustNewCaseWithMarker(
			&ptrByte,
			reflect.TypeOf(registrytest.PtrByte(nil)),
			true,
			core.TargetAny,
		),
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}

func TestIntConverter_Rune(t *testing.T) {
	conv := Integer()
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	equalFuncs := map[reflect.Type]func(got, want reflect.Value) bool{
		reflect.TypeFor[registrytest.Rune]():    equalRune,
		reflect.TypeFor[registrytest.PtrRune](): equalRune,
	}
	for typ, equalFunc := range equalFuncs {
		if err := tester.AddEqualFunc(typ, equalFunc); err != nil {
			t.Errorf("err occured: %s\n", err)
		}
	}
	r := markertest.NewMarker("rune", marker.STRING, reflect.ValueOf("r"))
	ptrRune := markertest.NewMarker("ptr.rune", marker.STRING, reflect.ValueOf("r"))
	cases := []convertertest.Case{
		tester.MustNewCaseWithMarker(&r, reflect.TypeOf(registrytest.Rune(0)), true, core.TargetAny),
		tester.MustNewCaseWithMarker(
			&ptrRune,
			reflect.TypeOf(registrytest.PtrRune(nil)),
			true,
			core.TargetAny,
		),
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}
