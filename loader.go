package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)
// TODO: I have to include eveything from types.* e.g. *types.Var oder
// *types.Signature etc.

type Loader interface {
	Load(patterns ...string) (map[string]*PackageInfo, error)
}

func NewLoader(conv Converter, cfg *packages.Config) Loader {
	l := &loader{
		conv:    conv,
		pkgInfo: &PackageInfo{},
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
	conv    Converter
	cfg     *packages.Config
	pkgInfo *PackageInfo
}

func (l *loader) Load(patterns ...string) (map[string]*PackageInfo, error) {
	info := make(map[string]*PackageInfo, len(patterns))
	pkgs, err := packages.Load(l.cfg, patterns...)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, errors.New("loaded packages are empty. Check that the defined patterns are matching any files")
	}
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, pkg.Errors[0]
		}
		for _, file := range pkg.Syntax {
			if file.Decls == nil {
				return nil, errors.New("no top-level declarations found")
			}
			if err := l.fileToInfo(pkg, file); err != nil {
				return nil, err
			}
			info[pkg.ID] = l.pkgInfo
			l.pkgInfo = &PackageInfo{}
		}
	}
	return info, nil
}

func (l *loader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes,
	}
}

func (l *loader) fileToInfo(pkg *packages.Package, file *ast.File) error {
	for _, decl := range file.Decls {
		funcDecl, isFuncDecl := decl.(*ast.FuncDecl)
		if isFuncDecl {
			l.funcDecl(funcDecl)
			continue
		}
		genDecl := decl.(*ast.GenDecl)
		switch genDecl.Tok {
		case token.CONST:
			l.constInfo(genDecl)
		case token.VAR:
			l.varInfo(genDecl)
		case token.IMPORT:
			l.importInfo(genDecl)
		case token.TYPE:
			l.typeDecl(pkg, genDecl)
		}
	}
	return nil
}

func (l *loader) funcDecl(fn *ast.FuncDecl) {
	if isMethod(fn) {
		l.methodInfo(fn)
		return
	}
	l.funcInfo(fn)
}

func (l *loader) typeDecl(pkg *packages.Package, gen *ast.GenDecl) {
	specs := convertSpecs[*ast.TypeSpec](gen.Specs)
	for _, spec := range specs {
		typ := pkg.TypesInfo.TypeOf(spec.Name)
		alias, isAlias := typ.(*types.Alias)
		if isAlias {
			l.aliasInfo(gen, alias, spec)
			continue
		}
		named := typ.(*types.Named).Underlying()
		pointer, isPointer := named.(*types.Pointer)
		if isPointer {
			l.typeInfo(gen, pointer.Elem(), spec)
			continue
		}
		strct, isStruct := named.(*types.Struct)
		if isStruct {
			l.structInfo(gen, strct, spec)
			continue
		}
		iface, isInterface := named.(*types.Interface)
		if isInterface {
			l.interfaceInfo(gen, iface, spec)
			continue
		}
	}
}

func (l *loader) funcInfo(fn *ast.FuncDecl) {
	info := &FuncInfo{
		doc:  fn.Doc.Text(),
		defs: Definitions{},
		Decl: fn,
	}
	l.pkgInfo.Funcs = append(l.pkgInfo.Funcs, info)
}

func (l *loader) methodInfo(fn *ast.FuncDecl) {
	info := &MethodInfo{
		doc:  fn.Doc.Text(),
		defs: Definitions{},
		Decl: fn,
	}
	l.pkgInfo.Methods = append(l.pkgInfo.Methods, info)
}

func (l *loader) constInfo(gen *ast.GenDecl) {
	specs := convertSpecs[*ast.ValueSpec](gen.Specs)
	for _, spec := range specs {
		doc := gen.Doc.Text() + spec.Doc.Text()
		info := &ConstInfo{
			doc:  doc,
			defs: Definitions{},
			Spec: spec,
			Decl: gen,
		}
		l.pkgInfo.Consts = append(l.pkgInfo.Consts, info)
	}
}

func (l *loader) varInfo(gen *ast.GenDecl) {
	specs := convertSpecs[*ast.ValueSpec](gen.Specs)
	for _, spec := range specs {
		doc := gen.Doc.Text() + spec.Doc.Text()
		info := &VarInfo{
			doc:  doc,
			defs: Definitions{},
			Spec: spec,
			Decl: gen,
		}
		l.pkgInfo.Vars = append(l.pkgInfo.Vars, info)
	}
}

func (l *loader) structInfo(gen *ast.GenDecl, strct *types.Struct, spec *ast.TypeSpec) {}

func (l *loader) interfaceInfo(gen *ast.GenDecl, iface *types.Interface, spec *ast.TypeSpec) {}

func (l *loader) importInfo(gen *ast.GenDecl) {}

func (l *loader) aliasInfo(gen *ast.GenDecl, alias *types.Alias, spec *ast.TypeSpec) {}

func (l *loader) typeInfo(gen *ast.GenDecl, typ types.Type, spec *ast.TypeSpec) {}
