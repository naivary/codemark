package loader

import (
	"fmt"
	"go/types"
	"testing"

	"golang.org/x/tools/go/packages"

	infov1 "github.com/naivary/codemark/api/info/v1"
	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/registry/registrytest"
)

// TODO: add unamed and named imports in the generated files
func TestLoaderLocal(t *testing.T) {
	tc, err := randLoaderTestCase()
	if err != nil {
		t.Errorf("err occured: %s", err)
	}
	cfg := &packages.Config{
		Dir: tc.Dir,
	}
	reg, err := registrytest.NewRegistry(registrytest.NewOptsSet())
	if err != nil {
		t.Errorf("err occured: %s", err)
	}
	mngr, err := converter.NewManager(reg)
	if err != nil {
		t.Errorf("err occured: %s", err)
	}
	l := New(mngr, cfg)
	projs, err := l.Load(".")
	if err != nil {
		t.Errorf("err occured: %s", err)
	}
	for _, proj := range projs {
		if err := isValid(tc, proj); err != nil {
			t.Errorf("err occured while reading %s: %s", tc.Dir, err)
		}
	}
}

func validateMarker(want []marker.Marker, got map[string][]any) error {
	for _, marker := range want {
		ident := marker.Ident
		_, found := got[ident]
		if !found {
			return fmt.Errorf("marker not found: %s\n", ident)
		}
	}
	return nil
}

func validate[T markers, V infov1.Info](typ string, want map[string]T, got map[types.Object]V) error {
	if len(got) != len(want) {
		return fmt.Errorf("quantity not equal for %s. got: %d; want: %d\n", typ, len(got), len(want))
	}
	for typ, info := range got {
		m := want[typ.Name()].markers()
		if err := validateMarker(m, info.Options()); err != nil {
			return err
		}
	}
	return nil
}

func isValid(tc loaderTestCase, proj *infov1.Information) error {
	// check struct
	if err := validate("structs", tc.Structs, proj.Structs); err != nil {
		return err
	}
	for typ, s := range proj.Structs {
		name := typ.Name()
		if err := validate("struct.fields", tc.Structs[name].Fields, s.Fields); err != nil {
			return err
		}
		if err := validate("struct.methods", tc.Structs[name].Methods, s.Methods); err != nil {
			return err
		}
	}
	// check iface
	if err := validate("interfaces", tc.Ifaces, proj.Ifaces); err != nil {
		return err
	}
	for typ, iface := range proj.Ifaces {
		name := typ.Name()
		if err := validate("interfaces.signature", tc.Ifaces[name].Signatures, iface.Signatures); err != nil {
			return err
		}
	}
	if err := validate("named", tc.Named, proj.Named); err != nil {
		return err
	}
	for typ, named := range proj.Named {
		name := typ.Name()
		if err := validate("named.methods", tc.Named[name].Methods, named.Methods); err != nil {
			return err
		}
	}
	// check rest
	if err := validate("consts", tc.Consts, proj.Consts); err != nil {
		return err
	}
	if err := validate("vars", tc.Vars, proj.Vars); err != nil {
		return err
	}
	if err := validate("imports", tc.Imports, proj.Imports); err != nil {
		return err
	}
	if err := validate("aliases", tc.Aliases, proj.Aliases); err != nil {
		return err
	}
	if err := validate("funcs", tc.Funcs, proj.Funcs); err != nil {
		return err
	}
	// check file because its a special case because of the missing types.Object
	if len(tc.Files) != len(proj.Files) {
		return fmt.Errorf("quantity not equal for files. got: %d; want: %d\n", len(proj.Files), len(tc.Files))
	}
	for filename, info := range proj.Files {
		markers := tc.Files[filename].markers()
		if err := validateMarker(markers, info.Opts); err != nil {
			return err
		}
	}
	return nil
}
