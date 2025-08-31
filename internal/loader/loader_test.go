package loader

import (
	"fmt"
	"go/types"
	"testing"

	"golang.org/x/tools/go/packages"

	infov1 "github.com/naivary/codemark/api/info/v1"
	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/internal/rand"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/registry/registrytest"
)

// TODO: add unamed and named imports in the generated files
func TestLoader_Local(t *testing.T) {
	proj := newProject()
	// addCustomDecls(proj)
	dir, err := proj.parse("tmpl/*")
	if err != nil {
		t.Errorf("unexpected error occured: %s", err)
	}
	cfg := &packages.Config{
		Dir: dir,
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
	information, err := l.Load(".")
	if err != nil {
		t.Errorf("err occured: %s", err)
	}
	for _, info := range information {
		if err := isValid(proj, info); err != nil {
			t.Errorf("err occured while reading %s: %s", dir, err)
		}
	}
}

func isValid(p *project, info *infov1.Information) error {
	// check struct
	if err := validate("structs", p.Structs, info.Structs); err != nil {
		return err
	}
	for typ, s := range info.Structs {
		name := typ.Name()
		if err := validate("struct.fields", p.Structs[name].Fields, s.Fields); err != nil {
			return err
		}
		if err := validate("struct.methods", p.Structs[name].Methods, s.Methods); err != nil {
			return err
		}
	}
	// check iface
	if err := validate("interfaces", p.Ifaces, info.Ifaces); err != nil {
		return err
	}
	for typ, iface := range info.Ifaces {
		name := typ.Name()
		if err := validate("interfaces.signature", p.Ifaces[name].Signatures, iface.Signatures); err != nil {
			return err
		}
	}
	if err := validate("named", p.Named, info.Named); err != nil {
		return err
	}
	for typ, named := range info.Named {
		name := typ.Name()
		if err := validate("named.methods", p.Named[name].Methods, named.Methods); err != nil {
			return err
		}
	}
	// check rest
	if err := validate("consts", p.Consts, info.Consts); err != nil {
		return err
	}
	if err := validate("vars", p.Vars, info.Vars); err != nil {
		return err
	}
	if err := validate("imports", p.Imports, info.Imports); err != nil {
		return err
	}
	if err := validate("aliases", p.Aliases, info.Aliases); err != nil {
		return err
	}
	if err := validate("funcs", p.Funcs, info.Funcs); err != nil {
		return err
	}
	// check file because its a special case because of the missing types.Object
	if len(p.Files) != len(info.Files) {
		return fmt.Errorf("quantity not equal for files. got: %d; want: %d\n", len(info.Files), len(p.Files))
	}
	for filename, info := range info.Files {
		markers := p.Files[filename].markers()
		if err := validateMarker(markers, info.Opts); err != nil {
			return err
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

func addCustomDecls(proj *project) {
	for _, imp := range importSpecial() {
		proj.Imports[imp.PackagePath] = imp
		proj.Vars[imp.Use.Ident] = imp.Use
	}
}

func importSpecial() []ImportDecl {
	return []ImportDecl{
		// aliased import
		{
			PackagePath: "fmt",
			Alias:       "format",
			Use: VarDecl{
				Ident:       rand.GoExpIdent(),
				Value:       "format.Println",
				IsImportUse: true,
			},
		},
	}
}
