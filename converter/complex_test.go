package converter

import (
	"testing"
)

func TestComplexConverter(t *testing.T) {
	tester, err := newConvTester(Complex())
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tests, err := tester.ValidTests()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range tests {
		tester.Run(t, tc, mngr)
	}
}
