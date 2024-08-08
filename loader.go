package main

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

// Loader is responsible for loading the specified
// files and their documentation
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
	for _, pkg := range pkgs {
		info := &Info{}
		consts, err := l.loadConsts(pkg)
		if err != nil {
			return nil, err
		}
		info.Consts = consts
		vars, err := l.loadVars(pkg)
		if err != nil {
			return nil, err
		}
		info.Vars = vars
		infos[pkg.ID] = info
	}
	return infos, nil
}

func (l *loader) loadConsts(pkg *packages.Package) ([]ConstInfo, error) {
	consts := make([]ConstInfo, 0, 0)
	for idn, obj := range pkg.TypesInfo.Defs {
		con, ok := obj.(*types.Const)
		if !ok {
			continue
		}
		info := ConstInfo{}
		info.Name = con.Name()
		info.Value = con.Val()
		info.Type = con.Type()
		info.Object = obj
		info.Ident = idn
		consts = append(consts, info)
	}
	return consts, nil
}

func (l *loader) loadVars(pkg *packages.Package) ([]VarInfo, error) {
	vars := make([]VarInfo, 0, 0)
	for _, obj := range pkg.TypesInfo.Defs {
		v, ok := obj.(*types.Var)
		if !ok {
			continue
		}
		if !l.isTopLevelVar(pkg, v) {
			continue
		}
		info := VarInfo{}
		info.Name = v.Name()
		info.Type = obj.Type()
		vars = append(vars, info)
	}
    fmt.Println(vars)
	return vars, nil
}

func (l *loader) isTopLevelVar(pkg *packages.Package, v *types.Var) bool {
	if v.IsField() {
		return false
	}
    // check if the variable is in a function declaration
	for _, file := range pkg.Syntax {
		pos := v.Pos()
		if !(file.FileStart <= pos && pos < file.FileEnd) {
			// not in this file
			continue
		}
		path, _ := astutil.PathEnclosingInterval(file, pos, pos)
		for _, n := range path {
			switch n.(type) {
			case *ast.FuncDecl:
				return false
			}
		}
	}

	return true
}

func (l *loader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypes,
	}
}
