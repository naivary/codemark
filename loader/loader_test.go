package loader

import (
	"fmt"
	"testing"

	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/maker"
	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/registry"
	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
	"golang.org/x/tools/go/packages"
)

func TestLoaderLocal(t *testing.T) {
	tc, err := RandLoaderTestCase()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	cfg := &packages.Config{
		Dir: tc.Dir,
	}
	reg, err := sdktesting.NewRegistry(registry.InMemory(), maker.New())
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	mngr, err := converter.NewManager(reg)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	l := New(mngr, cfg)
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

func isValid(tc LoaderTestCase, proj *sdk.Project) error {
	// check quantities
	if len(tc.Structs) != len(proj.Structs) {
		return fmt.Errorf("quantity of structs not equal. got: %d; want: %d", len(proj.Structs), len(tc.Structs))
	}
	for _, s := range proj.Structs {
		wantFields := tc.Structs[s.Spec.Name.Name].Fields
		if len(s.Fields) != len(wantFields) {
			return fmt.Errorf("quantity of fields not equal. got: %d; want: %d", len(s.Fields), len(wantFields))
		}
	}
	for _, s := range proj.Structs {
		wantMethods := tc.Structs[s.Spec.Name.Name].Methods
		if len(s.Methods) != len(wantMethods) {
			return fmt.Errorf("quantity of methods not equal. got: %d; want: %d", len(s.Methods), len(wantMethods))
		}
	}
	if len(tc.Funcs) != len(proj.Funcs) {
		return fmt.Errorf("quantity of funcs not equal. got: %d; want: %d", len(proj.Funcs), len(tc.Funcs))
	}
	if len(tc.Consts) != len(proj.Consts) {
		return fmt.Errorf("quantity of consts not equal. got: %d; want: %d", len(proj.Consts), len(tc.Consts))
	}
	if len(tc.Vars) != len(proj.Vars) {
		return fmt.Errorf("quantity of vars not equal. got: %d; want: %d", len(proj.Vars), len(tc.Vars))
	}
	if len(tc.Aliases) != len(proj.Aliases) {
		return fmt.Errorf("quantity of aliases not equal. got: %d; want: %d", len(proj.Aliases), len(tc.Aliases))
	}
	if len(tc.Ifaces) != len(proj.Ifaces) {
		return fmt.Errorf("quantity of interfaces not equal. got: %d; want: %d", len(proj.Ifaces), len(tc.Ifaces))
	}
	for _, iface := range proj.Ifaces {
		wantSig := tc.Ifaces[iface.Spec.Name.Name].Signatures
		if len(iface.Signatures) != len(wantSig) {
			return fmt.Errorf("quantity of signatures not equal. got: %d; want: %d", len(iface.Signatures), len(wantSig))
		}
	}
	if len(tc.Named) != len(proj.Named) {
		return fmt.Errorf("quantity of named not equal. got: %d; want: %d", len(proj.Named), len(tc.Named))
	}
	if len(tc.Imports) != len(proj.Imports) {
		return fmt.Errorf("quantity of imports not equal. got: %d; want: %d", len(proj.Imports), len(tc.Imports))
	}
	// check if the correct markers were loaded
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
