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

func (l *loader) loadConsts(pkg *packages.Package) ([]*ConstInfo, error) {
	consts := make([]*ConstInfo, 0, 0)
	for idn, obj := range pkg.TypesInfo.Defs {
		con, ok := obj.(*types.Const)
		if !ok {
			continue
		}
		info := &ConstInfo{}
		info.Name = con.Name()
		info.Value = con.Val()
		info.Type = con.Type()
		info.Object = obj
		info.Ident = idn
		consts = append(consts, info)
	}
	return consts, nil
}

func (l *loader) loadVars(pkg *packages.Package) ([]*VarInfo, error) {
	vars := make([]*VarInfo, 0, 0)
	for _, obj := range pkg.TypesInfo.Defs {
		v, ok := obj.(*types.Var)
		if !ok {
			continue
		}
		if !isVar(pkg, v) {
			continue
		}
		info := &VarInfo{}
		info.Name = v.Name()
		info.Type = obj.Type()
		vars = append(vars, info)
	}
	return vars, nil
}

func (l *loader) loadTypes(pkg *packages.Package) ([]*TypeInfo, error) {
	typs := make([]*TypeInfo, 0, 0)
	for _, obj := range pkg.TypesInfo.Defs {
		if obj == nil {
			continue
		}
		typeName, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}
		for _, file := range pkg.Syntax {
			pos := typeName.Pos()
			if !isInFile(file, pos) {
				continue
			}
			path, _ := astutil.PathEnclosingInterval(file, pos, pos)
			for _, n := range path {
				typeInfo := l.createTypeInfo(n, typeName)
				if typeInfo == nil {
					continue
				}
				l.createFieldInfos(pkg, typeInfo)
				typs = append(typs, typeInfo)
			}
		}
	}
	for _, t := range typs {
		fmt.Println(t)
	}
	return typs, nil
}

func (l *loader) createTypeInfo(node ast.Node, typeName *types.TypeName) *TypeInfo {
	decl, ok := node.(*ast.GenDecl)
	if !ok {
		return nil
	}
	if decl.Tok != token.TYPE {
		return nil
	}
	return newTypeInfo(typeName, decl)
}

func (l *loader) createFieldInfos(pkg *packages.Package, typeInfo *TypeInfo) {
	for _, spec := range typeInfo.GenDecl.Specs {
		v, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		typeInfo.Name = v.Name.Name
		if typeInfo.IsAlias {
			return
		}
		s, ok := v.Type.(*ast.StructType)
		if !ok {
			typeInfo.IsBasic = true
			return
		}
		typeInfo.IsStruct = true
		fieldInfos := make([]*FieldInfo, 0, 0)
		for _, field := range s.Fields.List {
			typ := pkg.TypesInfo.TypeOf(field.Type)
			info := newFieldInfo(field, typ)
			fieldInfos = append(fieldInfos, info...)
		}
		typeInfo.Fields = fieldInfos
	}
}

func (l *loader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypes,
	}
}

func isEmbedded(field *ast.Field) bool {
	return len(field.Names) == 0
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
