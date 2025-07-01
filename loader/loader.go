package loader

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/sdk"
	"golang.org/x/tools/go/packages"
)

var _ sdk.Loader = (*localLoader)(nil)

// TODO: add the current package to the struct so it doesnt have to be passed
// all the time
type localLoader struct {
	mngr sdk.ConverterManager
	cfg  *packages.Config
	proj *sdk.Project
}

// New Returns a new loader which can be used to read in go-packages.
func New(mngr sdk.ConverterManager, cfg *packages.Config) sdk.Loader {
	l := &localLoader{
		mngr: mngr,
		proj: sdk.NewProject(),
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

func (l *localLoader) Load(patterns ...string) (map[*packages.Package]*sdk.Project, error) {
	pkgs, err := packages.Load(l.cfg, patterns...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, errors.New("error occured after load")
	}
	if len(pkgs) == 0 {
		return nil, sdk.ErrPkgsEmpty
	}
	projs := make(map[*packages.Package]*sdk.Project, len(pkgs))
	for _, pkg := range pkgs {
		if err := l.exctractInfosFromPkg(pkg); err != nil {
			return nil, err
		}
		projs[pkg] = l.proj
		l.reset()
	}
	return projs, nil
}

func (l *localLoader) reset() {
	l.proj = sdk.NewProject()
}

func (l *localLoader) objectOf(pkg *packages.Package, ident *ast.Ident) (types.Object, error) {
	obj := pkg.TypesInfo.ObjectOf(ident)
	if obj == nil {
		return nil, fmt.Errorf("object not found: %v\n", ident)
	}
	return obj, nil
}

func (l *localLoader) exctractInfosFromPkg(pkg *packages.Package) error {
	for _, file := range pkg.Syntax {
		if err := l.extractInfosFromFile(pkg, file); err != nil {
			return err
		}
		if err := l.extractFileInfo(pkg, file); err != nil {
			return err
		}
	}
	return nil
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

func (l *localLoader) funcDecl(pkg *packages.Package, decl *ast.FuncDecl) error {
	if isMethod(decl) {
		return l.extractMethodInfo(pkg, decl)
	}
	return l.extractFuncInfo(pkg, decl)
}

func (l *localLoader) genDecl(pkg *packages.Package, decl *ast.GenDecl) error {
	var err error
	switch decl.Tok {
	case token.CONST:
		err = l.extractConstInfo(pkg, decl)
	case token.VAR:
		err = l.extractVarInfo(pkg, decl)
	case token.IMPORT:
		err = l.extractImportInfo(pkg, decl)
	case token.TYPE:
		err = l.typeDecl(pkg, decl)
	}
	return err
}

func (l *localLoader) typeDecl(pkg *packages.Package, decl *ast.GenDecl) error {
	specs := convertSpecs[*ast.TypeSpec](decl.Specs)
	for _, spec := range specs {
		var err error
		typ := pkg.TypesInfo.TypeOf(spec.Name)
		_, isAlias := typ.(*types.Alias)
		if isAlias {
			err = l.extractAliasInfo(pkg, decl, spec)
			if err != nil {
				return err
			}
			continue
		}
		switch typ.(*types.Named).Underlying().(type) {
		case *types.Struct:
			err = l.extractStructInfo(pkg, decl, spec)
		case *types.Interface:
			err = l.extractIfaceInfo(pkg, decl, spec)
		default:
			err = l.extractNamedTypeInfo(pkg, decl, spec)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *localLoader) extractMethodInfo(pkg *packages.Package, decl *ast.FuncDecl) error {
	doc := decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, target.METHOD)
	if err != nil {
		return err
	}
	info := sdk.FuncInfo{
		Pkg:  pkg,
		Decl: decl,
		Defs: defs,
	}
	rec := decl.Recv.List[0].Type
	recIdent, isIdent := rec.(*ast.Ident)
	if !isIdent {
		return fmt.Errorf("is not *ast.Ident: %v\n", rec)
	}
	recObj, err := l.objectOf(pkg, recIdent)
	if err != nil {
		return err
	}
	methObj, err := l.objectOf(pkg, decl.Name)
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

func (l *localLoader) extractFuncInfo(pkg *packages.Package, decl *ast.FuncDecl) error {
	doc := decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, target.FUNC)
	if err != nil {
		return err
	}
	obj, err := l.objectOf(pkg, decl.Name)
	if err != nil {
		return err
	}
	info := sdk.FuncInfo{
		Pkg:  pkg,
		Decl: decl,
		Defs: defs,
	}
	l.proj.Funcs[obj] = info
	return nil
}

func (l *localLoader) extractFileInfo(pkg *packages.Package, file *ast.File) error {
	doc := file.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, target.PKG)
	if err != nil {
		return err
	}
	info := sdk.FileInfo{
		File: file,
		Defs: defs,
	}

	filename := filepath.Base(pkg.Fset.Position(file.Package).Filename)
	l.proj.Files[filename] = info
	return nil
}

func (l *localLoader) extractImportInfo(pkg *packages.Package, decl *ast.GenDecl) error {
	specs := convertSpecs[*ast.ImportSpec](decl.Specs)
	for _, spec := range specs {
		doc := decl.Doc.Text() + spec.Doc.Text()
		defs, err := l.mngr.ParseDefs(doc, target.IMPORT)
		if err != nil {
			return err
		}
		info := sdk.ImportInfo{
			Pkg:  pkg,
			Spec: spec,
			Decl: decl,
			Defs: defs,
		}
		obj, err := l.objectOf(pkg, spec.Name)
		if err != nil {
			return err
		}
		l.proj.Imports[obj] = info
	}
	return nil
}

func (l *localLoader) extractVarInfo(pkg *packages.Package, decl *ast.GenDecl) error {
	specs := convertSpecs[*ast.ValueSpec](decl.Specs)
	for _, spec := range specs {
		doc := decl.Doc.Text() + spec.Doc.Text()
		for _, name := range spec.Names {
			defs, err := l.mngr.ParseDefs(doc, target.VAR)
			if err != nil {
				return err
			}
			info := sdk.VarInfo{
				Pkg:  pkg,
				Spec: spec,
				Decl: decl,
				Defs: defs,
			}
			obj, err := l.objectOf(pkg, name)
			if err != nil {
				return err
			}
			l.proj.Vars[obj] = info
		}
	}
	return nil
}

func (l *localLoader) extractConstInfo(pkg *packages.Package, decl *ast.GenDecl) error {
	specs := convertSpecs[*ast.ValueSpec](decl.Specs)
	for _, spec := range specs {
		doc := decl.Doc.Text() + spec.Doc.Text()
		for _, name := range spec.Names {
			defs, err := l.mngr.ParseDefs(doc, target.CONST)
			if err != nil {
				return err
			}
			info := sdk.ConstInfo{
				Pkg:  pkg,
				Spec: spec,
				Decl: decl,
				Defs: defs,
			}
			obj, err := l.objectOf(pkg, name)
			if err != nil {
				return err
			}
			l.proj.Consts[obj] = info
		}
	}
	return nil
}

func (l *localLoader) extractNamedTypeInfo(pkg *packages.Package, decl *ast.GenDecl, spec *ast.TypeSpec) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, target.NAMED)
	if err != nil {
		return err
	}
	info := sdk.NamedInfo{
		Pkg:     pkg,
		Spec:    spec,
		Decl:    decl,
		Defs:    defs,
		Methods: make(map[types.Object]sdk.FuncInfo),
	}
	obj, err := l.objectOf(pkg, spec.Name)
	if err != nil {
		return err
	}
	l.proj.Named[obj] = &info
	return nil
}

func (l *localLoader) extractIfaceInfo(pkg *packages.Package, decl *ast.GenDecl, spec *ast.TypeSpec) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, target.IFACE)
	if err != nil {
		return err
	}
	ifaceType := spec.Type.(*ast.InterfaceType)
	sigs, err := l.extractIfaceSignatureInfo(pkg, ifaceType)
	if err != nil {
		return err
	}
	info := sdk.IfaceInfo{
		Pkg:        pkg,
		Spec:       spec,
		Decl:       decl,
		Defs:       defs,
		Signatures: sigs,
	}
	obj, err := l.objectOf(pkg, spec.Name)
	if err != nil {
		return err
	}
	l.proj.Ifaces[obj] = info
	return nil
}

func (l *localLoader) extractIfaceSignatureInfo(pkg *packages.Package, spec *ast.InterfaceType) (map[types.Object]sdk.SignatureInfo, error) {
	infos := make(map[types.Object]sdk.SignatureInfo, spec.Methods.NumFields())
	for _, meth := range spec.Methods.List {
		doc := meth.Doc.Text()
		defs, err := l.mngr.ParseDefs(doc, target.IFACESIG)
		if err != nil {
			return nil, err
		}
		for _, name := range meth.Names {
			obj, err := l.objectOf(pkg, meth.Names[0])
			if err != nil {
				return nil, err
			}
			info := sdk.SignatureInfo{
				Idn:    name,
				Method: meth,
				Defs:   defs,
			}
			infos[obj] = info
		}
	}
	return infos, nil
}

func (l *localLoader) extractAliasInfo(pkg *packages.Package, decl *ast.GenDecl, spec *ast.TypeSpec) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, target.ALIAS)
	if err != nil {
		return err
	}
	info := sdk.AliasInfo{
		Pkg:  pkg,
		Decl: decl,
		Spec: spec,
		Defs: defs,
	}
	obj, err := l.objectOf(pkg, spec.Name)
	if err != nil {
		return err
	}
	l.proj.Aliases[obj] = info
	return nil
}

func (l *localLoader) extractStructInfo(pkg *packages.Package, decl *ast.GenDecl, spec *ast.TypeSpec) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, target.STRUCT)
	if err != nil {
		return err
	}
	obj, err := l.objectOf(pkg, spec.Name)
	if err != nil {
		return err
	}
	structType := spec.Type.(*ast.StructType)
	info := sdk.StructInfo{
		Defs:    defs,
		Spec:    spec,
		Decl:    decl,
		Pkg:     pkg,
		Fields:  make(map[types.Object]sdk.FieldInfo, structType.Fields.NumFields()),
		Methods: make(map[types.Object]sdk.FuncInfo, 0),
	}
	fieldInfos, err := l.extractFieldInfo(pkg, structType)
	if err != nil {
		return err
	}
	info.Fields = fieldInfos
	l.proj.Structs[obj] = &info
	return nil
}

func (l *localLoader) extractFieldInfo(pkg *packages.Package, spec *ast.StructType) (map[types.Object]sdk.FieldInfo, error) {
	infos := make(map[types.Object]sdk.FieldInfo, 0)
	for _, field := range spec.Fields.List {
		// embedded fields will be skipped
		if isEmbedded(field) {
			continue
		}
		doc := field.Doc.Text()
		defs, err := l.mngr.ParseDefs(doc, target.FIELD)
		if err != nil {
			return nil, err
		}
		for _, name := range field.Names {
			info := sdk.FieldInfo{
				Idn:   name,
				Field: field,
				Defs:  defs,
			}
			obj, err := l.objectOf(pkg, name)
			if err != nil {
				return nil, err
			}
			infos[obj] = info
		}
	}
	return infos, nil
}
