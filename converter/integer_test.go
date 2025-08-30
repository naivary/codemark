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

func TestIntConverter(t *testing.T) {
	conv := NewInteger()
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
	conv := NewInteger()
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	cases := []convertertest.Case{
		makeCase(Duration(0), markertest.New("duration", marker.STRING, "10h"), true),
		makeCase(PtrDuration(nil), markertest.New("ptr.duration", marker.STRING, "10h"), true),
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}

func TestIntConverter_Byte(t *testing.T) {
	conv := NewInteger()
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	cases := []convertertest.Case{
		makeCase(registrytest.Byte(0), markertest.New("byte", marker.STRING, "b"), true),
		makeCase(
			registrytest.PtrByte(nil),
			markertest.New("ptr.byte", marker.STRING, "b"),
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
	conv := NewInteger()
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	cases := []convertertest.Case{
		makeCase(registrytest.Rune(0), markertest.New("rune", marker.STRING, "r"), true),
		makeCase(
			registrytest.PtrRune(nil),
			markertest.New("ptr.rune", marker.STRING, "r"),
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
