package converter

import (
	"testing"
)

func TestAnyConverter(t *testing.T) {
	conv := Any()
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
