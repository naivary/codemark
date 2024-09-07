package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/packages"
)

type Loader interface {
	Load(patterns ...string) (map[string]*Infos, error)
}

func NewLoader(conv Converter, cfg *packages.Config) Loader {
	l := &loader{
		infos: NewInfos(),
		conv:  conv,
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
	conv  Converter
	cfg   *packages.Config
	infos *Infos
}

func (l *loader) Load(patterns ...string) (map[string]*Infos, error) {
	infos := make(map[string]*Infos, len(patterns))
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
			infos[pkg.ID] = l.infos
			l.infos = NewInfos()
		}
	}
	return infos, nil
}

func (l *loader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes,
	}
}

func (l *loader) fileToInfo(pkg *packages.Package, file *ast.File) error {
	return nil
}

/*func (l *loader) fileToInfo(pkg *packages.Package, file *ast.File) error {
	types := make([]*ast.GenDecl, 0, 0)
	for _, decl := range file.Decls {
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
	return l.loadDefs()
}*/
