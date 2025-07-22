package converter

import (
	"reflect"
	"testing"
	"time"

	coreapi "github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/converter/convertertest"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/typeutil"
)

type (
	Time    time.Time
	PtrTime *time.Time
)

func equalTime(got, want reflect.Value) bool {
	t := reflect.TypeOf(time.Time{})
	gotTime := typeutil.DerefValue(got).Convert(t).Interface().(time.Time)
	wantTime, err := time.Parse(time.RFC3339, want.String())
	if err != nil {
		return false
	}
	return gotTime.Equal(wantTime)
}

func customTimeCases(tester convertertest.Tester) []convertertest.Case {
	t := marker.New("codemark:testing:time", marker.STRING, reflect.ValueOf("2006-01-02T15:04:05Z"))
	ptrT := marker.New(
		"codemark:testing:ptr.time",
		marker.STRING,
		reflect.ValueOf("2006-01-02T15:04:05Z"),
	)

	return []convertertest.Case{
		tester.MustNewCaseWithMarker(&t, reflect.TypeFor[Time](), true, coreapi.TargetAny),
		tester.MustNewCaseWithMarker(&ptrT, reflect.TypeFor[PtrTime](), true, coreapi.TargetAny),
	}
}

func TestStringConverter(t *testing.T) {
	conv := String()
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

func TestStringConverter_Time(t *testing.T) {
	conv := String()
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	equalFuncs := map[reflect.Type]func(got, want reflect.Value) bool{
		reflect.TypeFor[Time]():    equalTime,
		reflect.TypeFor[PtrTime](): equalTime,
	}
	for to, equal := range equalFuncs {
		if err := tester.AddEqualFunc(to, equal); err != nil {
			t.Errorf("err occured: %s\n", err)
		}
	}
	for _, tc := range customTimeCases(tester) {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}
