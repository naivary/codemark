package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Loader interface {
	Load(files ...string) (map[string]*Info, error)
}

func NewLoader(cfg *packages.Config) Loader {
	l := &loader{
		info: NewInfo(),
	}
	if cfg == nil {
		l.cfg = l.defaultConfig()
	}
	l.cfg.ParseFile = func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
		return parser.ParseFile(fset, filename, src, parser.ParseComments)
	}
	return l
}

var _ Loader = (*loader)(nil)

type loader struct {
	cfg  *packages.Config
	info *Info
}

func (l *loader) Load(patterns ...string) (map[string]*Info, error) {
	infos := make(map[string]*Info, len(patterns))
	pkgs, err := packages.Load(l.cfg, patterns...)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, errors.New("empty packages")
	}
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, pkg.Errors[0]
		}
		for _, file := range pkg.Syntax {
			if err := l.fileToInfo(pkg, file); err != nil {
				return nil, err
			}
			infos[pkg.ID] = l.info
			l.info = NewInfo()
		}
	}
	return infos, nil
}

func (l *loader) fileToInfo(pkg *packages.Package, file *ast.File) error {
	if file.Decls == nil {
		return errors.New("no top-level declarations found")
	}
	types := make([]*ast.GenDecl, 0, 0)
	for _, decl := range file.Decls {
		// FuncDecl -> Method or Func
		funcDecl, isFuncDecl := decl.(*ast.FuncDecl)
		if isFuncDecl {
			l.funcDecl(funcDecl)
			continue
		}
		genDecl := decl.(*ast.GenDecl)
		switch genDecl.Tok {
		case token.CONST:
			l.constDecl(pkg, genDecl)
		case token.VAR:
			l.varDecl(pkg, genDecl)
		case token.IMPORT:
			l.importDecl(genDecl)
		case token.TYPE:
			types = append(types, genDecl)
		}
	}
	for _, typ := range types {
		l.typeDecl(pkg, typ)
	}
	return nil
}

func (l *loader) constDecl(pkg *packages.Package, decl *ast.GenDecl) {
	specs := convertSpecs[*ast.ValueSpec](decl.Specs)
	for _, spec := range specs {
		infos := NewConstInfo(spec, pkg)
		l.info.Consts = append(l.info.Consts, infos...)
	}
}

func (l *loader) varDecl(pkg *packages.Package, decl *ast.GenDecl) {
	specs := convertSpecs[*ast.ValueSpec](decl.Specs)
	for _, spec := range specs {
		infos := NewVarInfo(spec, pkg)
		l.info.Vars = append(l.info.Vars, infos...)
	}
}

func (l *loader) importDecl(decl *ast.GenDecl) {
	l.info.Imports = append(l.info.Imports, NewImportStmtInfo(decl))
}

func (l *loader) funcDecl(decl *ast.FuncDecl) {
	if decl.Recv != nil {
		l.methodDecl(decl)
		return
	}
	l.info.Funcs = append(l.info.Funcs, NewFuncInfo(decl))
}

func (l *loader) methodDecl(decl *ast.FuncDecl) {
	typ := decl.Recv.List[0].Type
	ident, isIdent := typ.(*ast.Ident)
	if isIdent {
		info := NewMethodInfo(decl, ident, nil)
		l.info.Methods[ident.Name] = append(l.info.Methods[ident.Name], info)
		return
	}
	pointer := typ.(*ast.StarExpr)
	ptrIdent := pointer.X.(*ast.Ident)
	info := NewMethodInfo(decl, ptrIdent, pointer)
	l.info.Methods[ptrIdent.Name] = append(l.info.Methods[ptrIdent.Name], info)
}

func (l *loader) typeDecl(pkg *packages.Package, decl *ast.GenDecl) {
	typeSpecs := convertSpecs[*ast.TypeSpec](decl.Specs)
	for _, typeSpec := range typeSpecs {
		typ := pkg.TypesInfo.TypeOf(typeSpec.Name)
		alias, isAlias := typ.(*types.Alias)
		if isAlias {
			info := NewAliasInfo(typeSpec, alias, decl)
			l.info.Aliases = append(l.info.Aliases, info)
			continue
		}

		named := typ.(*types.Named).Underlying()
		strc, isStruct := named.(*types.Struct)
		if isStruct {
			info := NewStructInfo(typeSpec, strc, decl, pkg)
			info.Methods = l.info.Methods[typeSpec.Name.Name]
			l.info.Structs = append(l.info.Structs, info)
			continue
		}
		iface, isIface := named.(*types.Interface)
		if isIface {
			info := NewInterfaceInfo(typeSpec, iface, decl)
			l.info.Interfaces = append(l.info.Interfaces, info)
			continue
		}

		basic, isBasic := named.(*types.Basic)
		if isBasic {
			info := NewBasicTypeInfo(typeSpec, basic, decl, nil)
			info.Methods = l.info.Methods[typeSpec.Name.Name]
			l.info.BasicTypes = append(l.info.BasicTypes, info)
			continue
		}

		ptr, isPtr := named.(*types.Pointer)
		if isPtr {
			basic := ptr.Elem().(*types.Basic)
			info := NewBasicTypeInfo(typeSpec, basic, decl, ptr)
			info.Methods = l.info.Methods[typeSpec.Name.Name]
			l.info.BasicTypes = append(l.info.BasicTypes, info)
		}
	}
}

func (l *loader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes,
	}
}
