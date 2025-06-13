package codemark

import (
	"testing"

	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func TestIntConverter(t *testing.T) {
	reg, err := sdktesting.NewDefsSet(NewInMemoryRegistry(), &DefinitionMarker{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	mngr, err := NewConvMngr(reg)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tester, err := sdktesting.NewConverterTester(&intConverter{}, nil)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tests, err := tester.Tests()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range tests {
		tester.Run(t, tc, mngr)
	}
}
