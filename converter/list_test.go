package converter

import (
	"testing"
)

func TestListConverter(t *testing.T) {
	tester, err := newConvTester(List(mngr))
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
