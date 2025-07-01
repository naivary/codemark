package converter

import (
	"testing"
)

func TestFloatConverter(t *testing.T) {
	tester, err := newConvTester(Float())
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}

	tests, err := tester.ValidTests()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range tests {
		tester.Run(t, tc)
	}
}
