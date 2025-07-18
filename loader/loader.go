package loader

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/naivary/codemark/api/core"
	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/converter"
	"golang.org/x/tools/go/packages"
)

var ErrPkgsEmpty = errors.New("loaded packages are empty. check that the defined patterns are matching any files")

type Loader interface {
	Load(patterns ...string) (map[*packages.Package]*loaderapi.Project, error)
}

var _ Loader = (*localLoader)(nil)

type localLoader struct {
	mngr *converter.Manager
	cfg  *packages.Config

	// proj is the current project being built
	proj *loaderapi.Project
	// pkg is the current package used to extract information from
	pkg *packages.Package
}

// New Returns a new loader which can be used to read in go-packages.
func New(mngr *converter.Manager, cfg *packages.Config) Loader {
	l := &localLoader{
		mngr: mngr,
		proj: loaderapi.NewProject(),
		cfg:  &packages.Config{},
	}
	if cfg != nil {
		l.cfg = cfg
	}
	l.cfg.Mode = packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes
	l.cfg.ParseFile = func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
		return parser.ParseFile(fset, filename, src, parser.ParseComments)
	}
	return l
}

func (l *localLoader) Load(patterns ...string) (map[*packages.Package]*loaderapi.Project, error) {
	pkgs, err := packages.Load(l.cfg, patterns...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, errors.New("error occured after load")
	}
	if len(pkgs) == 0 {
		return nil, ErrPkgsEmpty
	}
	projs := make(map[*packages.Package]*loaderapi.Project, len(pkgs))
	for _, pkg := range pkgs {
		l.pkg = pkg
		if err := l.exctractInfosFromPkg(); err != nil {
			return nil, err
		}
		projs[pkg] = l.proj
		l.reset()
	}
	return projs, nil
}

func (l *localLoader) reset() {
	l.proj = loaderapi.NewProject()
	l.pkg = nil
}

func (l *localLoader) objectOf(ident *ast.Ident) (types.Object, error) {
	obj := l.pkg.TypesInfo.ObjectOf(ident)
	if obj == nil {
		return nil, fmt.Errorf("object not found: %v\n", ident)
	}
	return obj, nil
}

func (l *localLoader) exctractInfosFromPkg() error {
	for _, file := range l.pkg.Syntax {
		if err := l.extractInfosFromFile(file); err != nil {
			return err
		}
		if err := l.extractFileInfo(file); err != nil {
			return err
		}
	}
	return nil
}

func (l *localLoader) extractInfosFromFile(file *ast.File) error {
	for _, decl := range file.Decls {
		var err error
		switch d := decl.(type) {
		case *ast.FuncDecl:
			err = l.funcDecl(d)
		case *ast.GenDecl:
			err = l.genDecl(d)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *localLoader) funcDecl(decl *ast.FuncDecl) error {
	if isMethod(decl) {
		return l.extractMethodInfo(decl)
	}
	return l.extractFuncInfo(decl)
}

func (l *localLoader) genDecl(decl *ast.GenDecl) error {
	var err error
	switch decl.Tok {
	case token.CONST:
		err = l.extractConstInfo(decl)
	case token.VAR:
		err = l.extractVarInfo(decl)
	case token.IMPORT:
		err = l.extractImportInfo(decl)
	case token.TYPE:
		err = l.typeDecl(decl)
	}
	return err
}

func (l *localLoader) typeDecl(decl *ast.GenDecl) error {
	specs := convertSpecs[*ast.TypeSpec](decl.Specs)
	for _, spec := range specs {
		var err error
		typ := l.pkg.TypesInfo.TypeOf(spec.Name)
		_, isAlias := typ.(*types.Alias)
		if isAlias {
			err = l.extractAliasInfo(decl, spec)
			if err != nil {
				return err
			}
			continue
		}
		switch typ.(*types.Named).Underlying().(type) {
		case *types.Struct:
			err = l.extractStructInfo(decl, spec)
		case *types.Interface:
			err = l.extractIfaceInfo(decl, spec)
		default:
			err = l.extractNamedTypeInfo(decl, spec)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *localLoader) extractMethodInfo(decl *ast.FuncDecl) error {
	doc := decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, core.TargetMethod)
	if err != nil {
		return err
	}
	info := loaderapi.FuncInfo{
		Decl: decl,
		Defs: defs,
	}
	rec := decl.Recv.List[0].Type
	recIdent, isIdent := rec.(*ast.Ident)
	if !isIdent {
		return fmt.Errorf("is not *ast.Ident: %v\n", rec)
	}
	recObj, err := l.objectOf(recIdent)
	if err != nil {
		return err
	}
	methObj, err := l.objectOf(decl.Name)
	if err != nil {
		return err
	}
	struc, isStruct := l.proj.Structs[recObj]
	if isStruct {
		struc.Methods[methObj] = info
		return nil
	}
	l.proj.Named[recObj].Methods[methObj] = info
	return nil
}

func (l *localLoader) extractFuncInfo(decl *ast.FuncDecl) error {
	doc := decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, core.TargetFunc)
	if err != nil {
		return err
	}
	obj, err := l.objectOf(decl.Name)
	if err != nil {
		return err
	}
	info := loaderapi.FuncInfo{
		Decl: decl,
		Defs: defs,
	}
	l.proj.Funcs[obj] = info
	return nil
}

func (l *localLoader) extractFileInfo(file *ast.File) error {
	doc := file.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, core.TargetPkg)
	if err != nil {
		return err
	}
	info := loaderapi.FileInfo{
		File: file,
		Defs: defs,
	}

	filename := filepath.Base(l.pkg.Fset.Position(file.Package).Filename)
	l.proj.Files[filename] = info
	return nil
}

func (l *localLoader) extractImportInfo(decl *ast.GenDecl) error {
	specs := convertSpecs[*ast.ImportSpec](decl.Specs)
	for _, spec := range specs {
		doc := decl.Doc.Text() + spec.Doc.Text()
		defs, err := l.mngr.ParseDefs(doc, core.TargetImport)
		if err != nil {
			return err
		}
		info := loaderapi.ImportInfo{
			Spec: spec,
			Decl: decl,
			Defs: defs,
		}
		obj, err := l.objectOf(spec.Name)
		if err != nil {
			return err
		}
		l.proj.Imports[obj] = info
	}
	return nil
}

func (l *localLoader) extractVarInfo(decl *ast.GenDecl) error {
	specs := convertSpecs[*ast.ValueSpec](decl.Specs)
	for _, spec := range specs {
		doc := decl.Doc.Text() + spec.Doc.Text()
		for _, name := range spec.Names {
			defs, err := l.mngr.ParseDefs(doc, core.TargetVar)
			if err != nil {
				return err
			}
			info := loaderapi.VarInfo{
				Spec: spec,
				Decl: decl,
				Defs: defs,
			}
			obj, err := l.objectOf(name)
			if err != nil {
				return err
			}
			l.proj.Vars[obj] = info
		}
	}
	return nil
}

func (l *localLoader) extractConstInfo(decl *ast.GenDecl) error {
	specs := convertSpecs[*ast.ValueSpec](decl.Specs)
	for _, spec := range specs {
		doc := decl.Doc.Text() + spec.Doc.Text()
		for _, name := range spec.Names {
			defs, err := l.mngr.ParseDefs(doc, core.TargetConst)
			if err != nil {
				return err
			}
			info := loaderapi.ConstInfo{
				Spec: spec,
				Decl: decl,
				Defs: defs,
			}
			obj, err := l.objectOf(name)
			if err != nil {
				return err
			}
			l.proj.Consts[obj] = info
		}
	}
	return nil
}

func (l *localLoader) extractNamedTypeInfo(decl *ast.GenDecl, spec *ast.TypeSpec) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, core.TargetNamed)
	if err != nil {
		return err
	}
	info := loaderapi.NamedInfo{
		Spec:    spec,
		Decl:    decl,
		Defs:    defs,
		Methods: make(map[types.Object]loaderapi.FuncInfo),
	}
	obj, err := l.objectOf(spec.Name)
	if err != nil {
		return err
	}
	l.proj.Named[obj] = &info
	return nil
}

func (l *localLoader) extractIfaceInfo(decl *ast.GenDecl, spec *ast.TypeSpec) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, core.TargetIface)
	if err != nil {
		return err
	}
	ifaceType := spec.Type.(*ast.InterfaceType)
	sigs, err := l.extractIfaceSignatureInfo(ifaceType)
	if err != nil {
		return err
	}
	info := loaderapi.IfaceInfo{
		Spec:       spec,
		Decl:       decl,
		Defs:       defs,
		Signatures: sigs,
	}
	obj, err := l.objectOf(spec.Name)
	if err != nil {
		return err
	}
	l.proj.Ifaces[obj] = info
	return nil
}

func (l *localLoader) extractIfaceSignatureInfo(spec *ast.InterfaceType) (map[types.Object]loaderapi.SignatureInfo, error) {
	infos := make(map[types.Object]loaderapi.SignatureInfo, spec.Methods.NumFields())
	for _, meth := range spec.Methods.List {
		doc := meth.Doc.Text()
		defs, err := l.mngr.ParseDefs(doc, core.TargetIfaceSig)
		if err != nil {
			return nil, err
		}
		for _, name := range meth.Names {
			obj, err := l.objectOf(meth.Names[0])
			if err != nil {
				return nil, err
			}
			info := loaderapi.SignatureInfo{
				Idn:    name,
				Method: meth,
				Defs:   defs,
			}
			infos[obj] = info
		}
	}
	return infos, nil
}

func (l *localLoader) extractAliasInfo(decl *ast.GenDecl, spec *ast.TypeSpec) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, core.TargetAlias)
	if err != nil {
		return err
	}
	info := loaderapi.AliasInfo{
		Decl: decl,
		Spec: spec,
		Defs: defs,
	}
	obj, err := l.objectOf(spec.Name)
	if err != nil {
		return err
	}
	l.proj.Aliases[obj] = info
	return nil
}

func (l *localLoader) extractStructInfo(decl *ast.GenDecl, spec *ast.TypeSpec) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, core.TargetStruct)
	if err != nil {
		return err
	}
	obj, err := l.objectOf(spec.Name)
	if err != nil {
		return err
	}
	structType := spec.Type.(*ast.StructType)
	info := loaderapi.StructInfo{
		Defs:    defs,
		Spec:    spec,
		Decl:    decl,
		Fields:  make(map[types.Object]loaderapi.FieldInfo, structType.Fields.NumFields()),
		Methods: make(map[types.Object]loaderapi.FuncInfo, 0),
	}
	fieldInfos, err := l.extractFieldInfo(structType)
	if err != nil {
		return err
	}
	info.Fields = fieldInfos
	l.proj.Structs[obj] = &info
	return nil
}

func (l *localLoader) extractFieldInfo(spec *ast.StructType) (map[types.Object]loaderapi.FieldInfo, error) {
	infos := make(map[types.Object]loaderapi.FieldInfo, 0)
	for _, field := range spec.Fields.List {
		// embedded fields will be skipped
		if isEmbedded(field) {
			continue
		}
		doc := field.Doc.Text()
		defs, err := l.mngr.ParseDefs(doc, core.TargetField)
		if err != nil {
			return nil, err
		}
		for _, name := range field.Names {
			info := loaderapi.FieldInfo{
				Idn:   name,
				Field: field,
				Defs:  defs,
			}
			obj, err := l.objectOf(name)
			if err != nil {
				return nil, err
			}
			infos[obj] = info
		}
	}
	return infos, nil
}
