package main

import (
	"errors"
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Loader interface {
	Load(files ...string) (map[string]*Info, error)
}

func NewLoader(cfg *packages.Config) Loader {
	l := &loader{}
	if cfg == nil {
		l.cfg = l.defaultConfig()
	}
	return l
}

var _ Loader = (*loader)(nil)

type loader struct {
	cfg *packages.Config

	vars           []*types.Var
	fields         []*types.Var
	aliases        []*types.Alias
	structs        []*types.Struct
	basicTypeNames []*types.TypeName
	funcs          []*types.Func
	methods        []*types.Func
}

// types und packages nutzen um die verschiedenen Expression reinzuladen
// docs für die expressions laden
// Marker für die Expression parsen und überprüfen ob die Marker auf diese
// Expression sein dürfen
// Info struct erstellen für die Expression e.g. FuncInfo oder ConstInfo etc.
// Eine Struct mit allen Info als result wiedergeben
func (l *loader) Load(paths ...string) (map[string]*Info, error) {
	infos := make(map[string]*Info, len(paths))
	pkgs, err := packages.Load(l.cfg, paths...)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, errors.New("empty packages")
	}
	for _, pkg := range pkgs {
		l.loadTypes(pkg)
	}
	return infos, nil
}

func (l *loader) loadTypes(pkg *packages.Package) {
	for _, obj := range pkg.TypesInfo.Defs {
		typ, isType := obj.(*types.TypeName)
		// for now ignore aliases
		if !isType || typ.IsAlias() {
			continue
		}
		named, isNamed := typ.Type().(*types.Named)
		if !isNamed {
			continue
		}
		basic, isBasic := named.Underlying().(*types.Basic)
		if isBasic {
			fmt.Println(basic)
		}
		struc, isStruct := named.Underlying().(*types.Struct)
        if isStruct {
            fmt.Println(struc)
        }
	}
}

func (l *loader) loadVars(pkg *packages.Package) {
	for _, obj := range pkg.TypesInfo.Defs {
		v, isVar := obj.(*types.Var)
		if !isVar {
			continue
		}
		fmt.Println(v)
	}
}

func (l *loader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes,
	}
}
