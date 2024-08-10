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
		info.Consts = l.loadConsts(pkg)
		info.Vars = l.loadVars(pkg)
		info.Types = l.loadTypes(pkg)
		info.Funcs = l.loadFuncs(pkg)
		infos[pkg.ID] = info
	}
	return infos, nil
}

func getNode[T any](syntax []*ast.File, pos token.Pos) []T {
	res := make([]T, 0, 0)
	for _, file := range syntax {
		if !isInFile(file, pos) {
			continue
		}
		path, _ := astutil.PathEnclosingInterval(file, pos, pos)
		for _, n := range path {
			v, ok := n.(T)
			if !ok {
				continue
			}
			res = append(res, v)
		}
	}
	return res
}

func (l *loader) loadFuncs(pkg *packages.Package) []*FuncInfo {
	infos := make([]*FuncInfo, 0, 0)
	for _, obj := range pkg.TypesInfo.Defs {
		fn, isFunc := obj.(*types.Func)
		if !isFunc {
			continue
		}
		decls := getNode[*ast.FuncDecl](pkg.Syntax, fn.Pos())
		infos = append(infos, l.createFuncInfo(decls)...)
	}
	return infos
}

func (l *loader) createFuncInfo(decls []*ast.FuncDecl) []*FuncInfo {
	funcInfos := make([]*FuncInfo, 0, len(decls))
	for _, decl := range decls {
		if isMethod(decl.Recv) {
			continue
		}
		fmt.Println(decl.Name)
		fmt.Println(decl.Doc.Text())
	}
	return funcInfos
}

func (l *loader) loadConsts(pkg *packages.Package) []*ConstInfo {
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
	return consts
}

func (l *loader) loadVars(pkg *packages.Package) []*VarInfo {
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
	return vars
}

func (l *loader) loadTypes(pkg *packages.Package) []*TypeInfo {
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
	return typs
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
			l.basicTypeInfo(v.Type, typeInfo)
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

func (l *loader) basicTypeInfo(expr ast.Expr, typeInfo *TypeInfo) {
	idn := expr.(*ast.Ident)
	typeInfo.IsBasic = true
	typeInfo.Ident = idn
}

func (l *loader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypes,
	}
}

func isEmbedded(field *ast.Field) bool {
	return len(field.Names) == 0
}

func isMethod(list *ast.FieldList) bool {
	if list == nil {
		return false
	}
	return true
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
