package converter

import (
	"testing"
	"time"

	"github.com/naivary/codemark/converter/convertertest"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/marker/markertest"
	"github.com/naivary/codemark/registry/registrytest"
)

type (
	Duration    time.Duration
	PtrDuration *time.Duration
)

// TODO: target cann raus beim convertertest.Tester
func TestIntConverter(t *testing.T) {
	conv := Integer()
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

func TestIntConverter_Duration(t *testing.T) {
	conv := Integer()
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	cases := []convertertest.Case{
		createCase(Duration(0), markertest.NewMarker("duration", marker.STRING, "10h"), true),
		createCase(PtrDuration(nil), markertest.NewMarker("ptr.duration", marker.STRING, "10h"), true),
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
	cases := []convertertest.Case{
		createCase(registrytest.Byte(0), markertest.NewMarker("byte", marker.STRING, "b"), true),
		createCase(
			registrytest.PtrByte(nil),
			markertest.NewMarker("ptr.byte", marker.STRING, "b"),
			true,
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
	cases := []convertertest.Case{
		createCase(registrytest.Rune(0), markertest.NewMarker("rune", marker.STRING, "r"), true),
		createCase(
			registrytest.PtrRune(nil),
			markertest.NewMarker("ptr.rune", marker.STRING, "r"),
			true,
		),
	}
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}
