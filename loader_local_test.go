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
	// TODO: Wrong order of want and got
	if len(tc.Structs) != len(proj.Structs) {
		return fmt.Errorf("quantity of structs not equal. got: %d; want: %d", len(tc.Structs), len(proj.Structs))
	}
	if len(tc.Funcs) != len(proj.Funcs) {
		return fmt.Errorf("quantity of funcs not equal. got: %d; want: %d", len(tc.Funcs), len(proj.Funcs))
	}
	if len(tc.Consts) != len(proj.Consts) {
		return fmt.Errorf("quantity of consts not equal. got: %d; want: %d", len(tc.Consts), len(proj.Consts))
	}
	if len(tc.Vars) != len(proj.Vars) {
		return fmt.Errorf("quantity of vars not equal. got: %d; want: %d", len(tc.Vars), len(proj.Vars))
	}
	if len(tc.Aliases) != len(proj.Aliases) {
		return fmt.Errorf("quantity of aliases not equal. got: %d; want: %d", len(tc.Aliases), len(proj.Aliases))
	}
	if len(tc.Ifaces) != len(proj.Ifaces) {
		return fmt.Errorf("quantity of interfaces not equal. got: %d; want: %d", len(tc.Ifaces), len(proj.Ifaces))
	}
	if len(tc.Named) != len(proj.Named) {
		return fmt.Errorf("quantity of named not equal. got: %d; want: %d", len(tc.Named), len(proj.Named))
	}
	if len(tc.Imports) != len(proj.Imports) {
		return fmt.Errorf("quantity of imports not equal. got: %d; want: %d", len(tc.Imports), len(proj.Imports))
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
	for typ, v := range proj.Vars {
		name := typ.Name()
		want := tc.Vars[name]
		if err := isMarkerMatching(want.Markers, v.Defs); err != nil {
			return err
		}
	}
	for typ, a := range proj.Aliases {
		name := typ.Name()
		want := tc.Aliases[name]
		if err := isMarkerMatching(want.Markers, a.Defs); err != nil {
			return err
		}
	}
	for typ, i := range proj.Ifaces {
		name := typ.Name()
		want := tc.Ifaces[name]
		if err := isMarkerMatching(want.Markers, i.Defs); err != nil {
			return err
		}
	}
	for typ, n := range proj.Named {
		name := typ.Name()
		want := tc.Named[name]
		if err := isMarkerMatching(want.Markers, n.Defs); err != nil {
			return err
		}
	}
	for typ, im := range proj.Imports {
		name := typ.Name()
		want := tc.Imports[name]
		if err := isMarkerMatching(want.Markers, im.Defs); err != nil {
			return err
		}
	}
	for typ, pkg := range proj.Pkgs {
		want := tc.Imports[typ.Name]
		if err := isMarkerMatching(want.Markers, pkg.Defs); err != nil {
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
