package converter

import (
	"testing"
)

// ADD test for rune and byte e.g. from "s" to byte
func TestIntConverter(t *testing.T) {
	tester, err := newConvTester(&intConverter{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tests, err := tester.ValidTests()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}
