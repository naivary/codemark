package main

import (
	"fmt"
	"go/ast"
	"go/token"
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
		types, err := l.loadTypes(pkg)
		if err != nil {
			return nil, err
		}
		info.Consts = consts
		info.Vars = vars
		info.Types = types
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
		if !isVar(pkg, v) {
			continue
		}
		info := VarInfo{}
		info.Name = v.Name()
		info.Type = obj.Type()
		vars = append(vars, info)
	}
	return vars, nil
}

func (l *loader) loadTypes(pkg *packages.Package) ([]TypeInfo, error) {
	typs := make([]TypeInfo, 0, 0)
	for _, obj := range pkg.TypesInfo.Defs {
		if obj == nil {
			continue
		}
		v, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}
		if v.IsAlias() {
			continue
		}
		for _, file := range pkg.Syntax {
			pos := v.Pos()
			if !isInFile(file, pos) {
				continue
			}
			path, _ := astutil.PathEnclosingInterval(file, pos, pos)
			for _, n := range path {
                _ = n
			}
		}
	}
	return typs, nil
}

func (l *loader) createTypeInfo(node ast.Node) *TypeInfo {
	info := &TypeInfo{}
	decl, ok := node.(*ast.GenDecl)
	if !ok {
		return nil
	}
	if decl.Tok != token.TYPE {
		return nil
	}
	info.Doc = decl.Doc.Text()
	info.GenDecl = decl
	return info
}

func (l *loader) createFieldInfos(specs []ast.Spec) []FieldInfo {
	for _, spec := range specs {
		v, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		t, ok := v.Type.(*ast.StructType)
		fmt.Println(t)
	}
	return nil
}

func (l *loader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypes,
	}
}

func isInFile(file *ast.File, pos token.Pos) bool {
	return file.FileStart <= pos && pos < file.FileEnd
}

func isVar(pkg *packages.Package, v *types.Var) bool {
	if v.IsField() {
		return false
	}
	// check if the variable is in a function declaration
	for _, file := range pkg.Syntax {
		pos := v.Pos()
		if !isInFile(file, pos) {
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
