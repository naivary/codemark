package codemark

import (
	"fmt"
	"testing"

	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
	"golang.org/x/tools/go/packages"
)

func TestLoaderLocal(t *testing.T) {
	tc, err := sdktesting.NewGoFiles()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	cfg := &packages.Config{
		Dir: tc.Dir,
	}
	reg, err := sdktesting.NewDefsSet(NewInMemoryRegistry(), &DefinitionMarker{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	mngr, err := NewConvMngr(reg)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	l := NewLocalLoader(mngr, cfg)
	projs, err := l.Load(".")
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, proj := range projs {
		if err := isValid(tc, proj); err != nil {
			t.Errorf("err occured: %s\n", err)
		}
	}
}

// TODO: what should be the logic to verify the corretnes of the loader?
// check the number of got and wanted informations about types
// did it load the correct marker?
func isValid(tc sdktesting.LoaderTestCase, proj *sdk.Project) error {
	// check quantities
	if len(tc.Structs) != len(proj.Structs) {
		return fmt.Errorf("quantity of structs not equal. got: %d; want: %d", len(tc.Structs), len(proj.Structs))
	}
	if len(tc.Funcs) != len(proj.Funcs) {
		return fmt.Errorf("quantity of funcs not equal. got: %d; want: %d", len(tc.Funcs), len(proj.Funcs))
	}

	return nil
}
