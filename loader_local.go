package codemark

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/naivary/codemark/sdk"
	"golang.org/x/tools/go/packages"
)

var _ sdk.Loader = (*localLoader)(nil)

type localLoader struct {
	mngr *ConverterManager
	cfg  *packages.Config
	proj *sdk.Project
}

func NewLocalLoader(mngr *ConverterManager, cfg *packages.Config) sdk.Loader {
	l := &localLoader{
		mngr: mngr,
		proj: &sdk.Project{},
	}
	if cfg == nil {
		l.cfg = l.defaultConfig()
	}
	l.cfg.ParseFile = func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
		return parser.ParseFile(fset, filename, src, parser.ParseComments)
	}
	return l
}

func (l *localLoader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes,
	}
}

func (l *localLoader) Load(patterns ...string) (*sdk.Project, error) {
	pkgs, err := packages.Load(l.cfg, patterns...)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, sdk.ErrPkgsEmpty
	}
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, pkg.Errors[0]
		}
		for _, file := range pkg.Syntax {
			if err := l.extractInfosFromFile(pkg, file); err != nil {
				return nil, err
			}
		}
	}
	return l.proj, nil
}

func (l *localLoader) extractInfosFromFile(pkg *packages.Package, file *ast.File) error {
	for _, decl := range file.Decls {
		var err error
		switch d := decl.(type) {
		case *ast.FuncDecl:
			err = l.funcDecl(pkg, d)
		case *ast.GenDecl:
			err = l.genDecl(pkg, d)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *localLoader) funcDecl(_ *packages.Package, decl *ast.FuncDecl) error {
	if isMethod(decl) {
		return nil
	}
	return nil
}

func (l *localLoader) genDecl(pkg *packages.Package, decl *ast.GenDecl) error {
	var err error
	switch decl.Tok {
	case token.CONST:
		//err = l.constInfo(pkg, decl)
	case token.VAR:
		//err = l.varInfo(pkg, decl)
	case token.IMPORT:
		//err = l.importInfo(pkg, decl)
	case token.TYPE:
		err = l.typeDecl(pkg, decl)
	}
	if err != nil {
		return err
	}
	return nil
}

func (l *localLoader) typeDecl(pkg *packages.Package, decl *ast.GenDecl) error {
	specs := convertSpecs[*ast.TypeSpec](decl.Specs)
	for _, spec := range specs {
		var err error
		typ := pkg.TypesInfo.TypeOf(spec.Name)
		switch typ.(*types.Named).Underlying().(type) {
		// case *types.Alias:
		case *types.Struct:
			err = l.extractStructInfo(pkg, decl, spec)
		//case *types.Interface:
		default:
			// everything else beside interface, struct and alias
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *localLoader) extractStructInfo(pkg *packages.Package, decl *ast.GenDecl, spec *ast.TypeSpec) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, sdk.TargetStruct)
	if err != nil {
		return err
	}
	structType := spec.Type.(*ast.StructType)
	info := sdk.StructInfo{
		Defs:   defs,
		Spec:   spec,
		Decl:   decl,
		Pkg:    pkg,
		Fields: make([]sdk.FieldInfo, 0, structType.Fields.NumFields()),
	}
	fieldInfos, err := l.extractFieldInfo(structType)
	if err != nil {
		return err
	}
	info.Fields = fieldInfos
	l.proj.Structs = append(l.proj.Structs, info)
	return nil
}

func (l *localLoader) extractFieldInfo(spec *ast.StructType) ([]sdk.FieldInfo, error) {
	infos := make([]sdk.FieldInfo, 0, 0)
	for _, field := range spec.Fields.List {
		doc := field.Doc.Text()
		defs, err := l.mngr.ParseDefs(doc, sdk.TargetField)
		if err != nil {
			return nil, err
		}
		for _, name := range field.Names {
			info := sdk.FieldInfo{
				Idn:   name,
				Field: field,
				Defs:  defs,
			}
			infos = append(infos, info)
		}
	}
	return infos, nil
}
