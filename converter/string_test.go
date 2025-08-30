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
	conv := NewString()
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
	conv := NewString()
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	cases := []convertertest.Case{
		makeCase(
			Time(time.Time{}),
			markertest.New("time", marker.STRING, "2006-01-02T15:04:05Z"),
			true,
		),
		makeCase(
			PtrTime(nil),
			markertest.New("ptr.time", marker.STRING, "2006-01-02T15:04:05Z"),
			true,
		),
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}
