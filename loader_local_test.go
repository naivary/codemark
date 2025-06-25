package codemark

import (
	"fmt"
	"testing"

	"github.com/naivary/codemark/parser"
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
// did it load the correct marker? (names only because values are not the
// responsiblities of the loader but the converter)
func isValid(tc sdktesting.LoaderTestCase, proj *sdk.Project) error {
	// check quantities
	if len(tc.Structs) != len(proj.Structs) {
		return fmt.Errorf("quantity of structs not equal. got: %d; want: %d", len(tc.Structs), len(proj.Structs))
	}
	if len(tc.Funcs) != len(proj.Funcs) {
		return fmt.Errorf("quantity of funcs not equal. got: %d; want: %d", len(tc.Funcs), len(proj.Funcs))
	}
	if len(tc.Consts) != len(proj.Consts) {
		return fmt.Errorf("quantity of const not equal. got: %d; want: %d", len(tc.Consts), len(proj.Consts))
	}

	// check marker existence
	for typ, fn := range proj.Funcs {
		name := typ.Name()
		want := tc.Funcs[name]
		if err := isMarkerMatching(want.Markers, fn.Defs); err != nil {
			return err
		}
	}
	for typ, stc := range proj.Structs {
		name := typ.Name()
		want := tc.Structs[name]
		if err := isMarkerMatching(want.Markers, stc.Defs); err != nil {
			return err
		}
	}
	for typ, c := range proj.Consts {
		name := typ.Name()
		want := tc.Consts[name]
		if err := isMarkerMatching(want.Markers, c.Defs); err != nil {
			return err
		}
	}

	return nil
}

func isMarkerMatching(want []parser.Marker, got map[string][]any) error {
	for _, marker := range want {
		ident := marker.Ident()
		_, found := got[ident]
		if !found {
			return fmt.Errorf("marker not found: %s\n", ident)
		}
	}
	return nil
}
