package converter

import (
	"testing"
	"time"

	"github.com/naivary/codemark/converter/convertertest"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/marker/markertest"
)

type (
	Time    time.Time
	PtrTime *time.Time
)

func TestStringConverter(t *testing.T) {
	conv := String()
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	cases, err := validCasesFor(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range cases {
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
	cases := []convertertest.Case{
		createCase(
			Time(time.Time{}),
			markertest.NewMarker("time", marker.STRING, "2006-01-02T15:04:05Z"),
			true,
		),
		createCase(
			PtrTime(nil),
			markertest.NewMarker("ptr.time", marker.STRING, "2006-01-02T15:04:05Z"),
			true,
		),
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}
