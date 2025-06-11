package codemark

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
	"golang.org/x/tools/go/packages"
)

var _ sdk.Loader = (*localLoader)(nil)

type localLoader struct {
	mngr    *ConverterManager
	cfg     *packages.Config
	proj    *sdk.Project
	methods map[string][]sdk.FuncInfo
}

func NewLocalLoader(mngr *ConverterManager, cfg *packages.Config) sdk.Loader {
	l := &localLoader{
		mngr:    mngr,
		proj:    &sdk.Project{},
		methods: make(map[string][]sdk.FuncInfo),
	}
	if cfg == nil {
		l.cfg = l.defaultConfig()
	}
	l.cfg.ParseFile = func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
		return parser.ParseFile(fset, filename, src, parser.ParseComments)
	}
	return l
}

func (l *localLoader) Load(patterns ...string) ([]*sdk.Project, error) {
	pkgs, err := packages.Load(l.cfg, patterns...)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, sdk.ErrPkgsEmpty
	}
	projs := make([]*sdk.Project, 0, len(pkgs))
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, pkg.Errors[0]
		}
		for _, file := range pkg.Syntax {
			if err := l.extractInfosFromFile(pkg, file); err != nil {
				return nil, err
			}
			if err := l.extractPkgInfo(pkg, file); err != nil {
				return nil, err
			}
			for _, struc := range l.proj.Structs {
				name := struc.Spec.Name.Name
				struc.Methods = l.methods[name]
			}
			for _, named := range l.proj.Named {
				name := named.Spec.Name.Name
				named.Methods = l.methods[name]
			}
		}
		projs = append(projs, l.proj)
		// reset
		l.proj = &sdk.Project{}
		l.methods = make(map[string][]sdk.FuncInfo)
	}
	return projs, nil
}

func (l *localLoader) defaultConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes,
	}
}

func (l *localLoader) addMethod(info sdk.FuncInfo) {
	expr := info.Decl.Recv.List[0].Type
	name := sdkutil.ExprName(expr)
	infos, ok := l.methods[name]
	if ok {
		l.methods[name] = append(infos, info)
		return
	}
	l.methods[name] = []sdk.FuncInfo{info}
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
	if sdkutil.IsMethod(decl) {
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
	specs := sdkutil.ConvertSpecs[*ast.TypeSpec](decl.Specs)
	for _, spec := range specs {
		var err error
		typ := pkg.TypesInfo.TypeOf(spec.Name)
		switch typ.(*types.Named).Underlying().(type) {
		case *types.Alias:
			err = l.extractAliasInfo(pkg, decl, spec)
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
	defs, err := l.mngr.ParseDefs(doc, sdk.TargetType)
	if err != nil {
		return err
	}
	info := sdk.FuncInfo{
		Pkg:  pkg,
		Decl: decl,
		Defs: defs,
	}
	l.addMethod(info)
	return nil
}

func (l *localLoader) extractFuncInfo(pkg *packages.Package, decl *ast.FuncDecl) error {
	doc := decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, sdk.TargetType)
	if err != nil {
		return err
	}
	info := sdk.FuncInfo{
		Pkg:  pkg,
		Decl: decl,
		Defs: defs,
	}
	l.proj.Funcs = append(l.proj.Funcs, info)
	return nil
}

func (l *localLoader) extractPkgInfo(pkg *packages.Package, file *ast.File) error {
	doc := file.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, sdk.TargetType)
	if err != nil {
		return err
	}
	info := sdk.PkgInfo{
		Pkg:  pkg,
		File: file,
		Defs: defs,
	}
	l.proj.Pkgs = append(l.proj.Pkgs, info)
	return nil
}

func (l *localLoader) extractImportInfo(pkg *packages.Package, decl *ast.GenDecl) error {
	specs := sdkutil.ConvertSpecs[*ast.ImportSpec](decl.Specs)
	for _, spec := range specs {
		doc := decl.Doc.Text() + spec.Doc.Text()
		defs, err := l.mngr.ParseDefs(doc, sdk.TargetImport)
		if err != nil {
			return err
		}
		info := sdk.ImportInfo{
			Pkg:  pkg,
			Spec: spec,
			Decl: decl,
			Defs: defs,
		}
		l.proj.Imports = append(l.proj.Imports, info)
	}
	return nil
}

func (l *localLoader) extractVarInfo(pkg *packages.Package, decl *ast.GenDecl) error {
	specs := sdkutil.ConvertSpecs[*ast.ValueSpec](decl.Specs)
	for _, spec := range specs {
		doc := decl.Doc.Text() + spec.Doc.Text()
		defs, err := l.mngr.ParseDefs(doc, sdk.TargetVar)
		if err != nil {
			return err
		}
		info := sdk.VarInfo{
			Pkg:  pkg,
			Spec: spec,
			Decl: decl,
			Defs: defs,
		}
		l.proj.Vars = append(l.proj.Vars, info)
	}
	return nil
}

func (l *localLoader) extractConstInfo(pkg *packages.Package, decl *ast.GenDecl) error {
	specs := sdkutil.ConvertSpecs[*ast.ValueSpec](decl.Specs)
	for _, spec := range specs {
		doc := decl.Doc.Text() + spec.Doc.Text()
		defs, err := l.mngr.ParseDefs(doc, sdk.TargetConst)
		if err != nil {
			return err
		}
		info := sdk.ConstInfo{
			Pkg:  pkg,
			Spec: spec,
			Decl: decl,
			Defs: defs,
		}
		l.proj.Consts = append(l.proj.Consts, info)
	}
	return nil
}

func (l *localLoader) extractNamedTypeInfo(pkg *packages.Package, decl *ast.GenDecl, spec *ast.TypeSpec) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, sdk.TargetType)
	if err != nil {
		return err
	}
	info := sdk.NamedInfo{
		Pkg:  pkg,
		Spec: spec,
		Decl: decl,
		Defs: defs,
	}
	l.proj.Named = append(l.proj.Named, &info)
	return nil
}

func (l *localLoader) extractIfaceInfo(pkg *packages.Package, decl *ast.GenDecl, spec *ast.TypeSpec) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, sdk.TargetInterface)
	if err != nil {
		return err
	}
	ifaceType := spec.Type.(*ast.InterfaceType)
	sigs, err := l.extractIfaceSignatureInfo(ifaceType)
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
	l.proj.Ifaces = append(l.proj.Ifaces, info)
	return nil
}

func (l *localLoader) extractIfaceSignatureInfo(spec *ast.InterfaceType) ([]sdk.SignatureInfo, error) {
	infos := make([]sdk.SignatureInfo, 0, spec.Methods.NumFields())
	for _, meth := range spec.Methods.List {
		doc := meth.Doc.Text()
		defs, err := l.mngr.ParseDefs(doc, sdk.TargetInterfaceSignature)
		if err != nil {
			return nil, err
		}
		for _, name := range meth.Names {
			info := sdk.SignatureInfo{
				Idn:    name,
				Method: meth,
				Defs:   defs,
			}
			infos = append(infos, info)
		}
	}
	return infos, nil
}

func (l *localLoader) extractAliasInfo(pkg *packages.Package, decl *ast.GenDecl, spec *ast.TypeSpec) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	defs, err := l.mngr.ParseDefs(doc, sdk.TargetAlias)
	if err != nil {
		return err
	}
	info := sdk.AliasInfo{
		Pkg:  pkg,
		Decl: decl,
		Spec: spec,
		Defs: defs,
	}
	l.proj.Aliases = append(l.proj.Aliases, info)
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
	l.proj.Structs = append(l.proj.Structs, &info)
	return nil
}

func (l *localLoader) extractFieldInfo(spec *ast.StructType) ([]sdk.FieldInfo, error) {
	infos := make([]sdk.FieldInfo, 0, 0)
	for _, field := range spec.Fields.List {
		// embedded fields will be skipped
		if sdkutil.IsEmbedded(field) {
			continue
		}
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
