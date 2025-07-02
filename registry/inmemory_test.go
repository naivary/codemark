package registry

import (
	"testing"

	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func TestInMemory(t *testing.T) {
	reg := InMemory()
	tester := sdktesting.NewRegTester(reg)
	defs := sdktesting.NewDefSet()
	tests := tester.ValidTests(defs...)
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}
