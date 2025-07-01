package loader

import (
	"fmt"
	"go/types"
	"testing"

	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/parser"
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
	reg, err := sdktesting.NewRegistry(nil)
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
			t.Errorf("err occured while reading %s: %s\n", tc.Dir, err)
		}
	}
}

func validateMarker(want []parser.Marker, got map[string][]any) error {
	for _, marker := range want {
		ident := marker.Ident()
		_, found := got[ident]
		if !found {
			return fmt.Errorf("marker not found: %s\n", ident)
		}
	}
	return nil
}

func validate[T markers, V defs](want map[string]T, got map[types.Object]V) error {
	if len(got) != len(want) {
		return fmt.Errorf("quantity not equal. got: %d; want: %d\n", len(got), len(want))
	}
	for typ, info := range got {
		m := want[typ.Name()].markers()
		if err := validateMarker(m, info.Definitions()); err != nil {
			return err
		}
	}
	return nil
}

func isValid(tc LoaderTestCase, proj *sdk.Project) error {
	// check struct
	if err := validate(tc.Structs, proj.Structs); err != nil {
		return err
	}
	for typ, s := range proj.Structs {
		name := typ.Name()
		if err := validate(tc.Structs[name].Fields, s.Fields); err != nil {
			return err
		}
		if err := validate(tc.Structs[name].Methods, s.Methods); err != nil {
			return err
		}
	}
	// check iface
	if err := validate(tc.Ifaces, proj.Ifaces); err != nil {
		return err
	}
	for typ, iface := range proj.Ifaces {
		name := typ.Name()
		if err := validate(tc.Ifaces[name].Signatures, iface.Signatures); err != nil {
			return err
		}
	}
	if err := validate(tc.Named, proj.Named); err != nil {
		return err
	}
	for typ, named := range proj.Named {
		name := typ.Name()
		if err := validate(tc.Named[name].Methods, named.Methods); err != nil {
			return err
		}
	}
	// check rest
	if err := validate(tc.Consts, proj.Consts); err != nil {
		return err
	}
	if err := validate(tc.Vars, proj.Vars); err != nil {
		return err
	}
	if err := validate(tc.Imports, proj.Imports); err != nil {
		return err
	}
	if err := validate(tc.Aliases, proj.Aliases); err != nil {
		return err
	}
	if err := validate(tc.Funcs, proj.Funcs); err != nil {
		return err
	}
	// check filee bcause its a special case because of the missing types.Object
	if len(tc.Files) != len(proj.Files) {
		return fmt.Errorf("quantity not equal. got: %d; want: %d\n", len(proj.Files), len(tc.Files))
	}
	for filename, info := range proj.Files {
		markers := tc.Files[filename].markers()
		if err := validateMarker(markers, info.Defs); err != nil {
			return err
		}
	}
	return nil
}
