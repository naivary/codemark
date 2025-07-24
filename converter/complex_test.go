package converter

import (
	"testing"

	"github.com/naivary/codemark/converter/convertertest"
)

func TestComplexConverter(t *testing.T) {
	conv := NewComplex()
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
