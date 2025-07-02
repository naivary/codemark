package registry

import (
	"testing"

	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func TestInMemory(t *testing.T) {
	tester := sdktesting.NewRegTester(InMemory())
	defs := sdktesting.NewDefSet()
	tests := tester.ValidTests(defs...)
	for _, tc := range tests {
		tester.Run(t, tc)
	}
}
