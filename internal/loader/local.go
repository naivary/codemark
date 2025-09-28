package loader

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"

	"golang.org/x/tools/go/packages"

	infov1 "github.com/naivary/codemark/api/info/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
	"github.com/naivary/codemark/converter"
)

type parseMarkers = func(doc string, target optionv1.Target) (infov1.Options, error)

var _ Loader = (*loader)(nil)

type loader struct {
	mngr *converter.Manager

	cfg *packages.Config
}

// New Returns a new loader which can be used to read in go-packages.
func New(mngr *converter.Manager, cfg *packages.Config) Loader {
	if cfg == nil {
		cfg = &packages.Config{}
	}
	l := &loader{
		mngr: mngr,
		cfg: &packages.Config{
			Mode: packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes | packages.NeedImports | packages.NeedName,
			ParseFile: func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
				return parser.ParseFile(fset, filename, src, parser.ParseComments)
			},
			Dir:        cfg.Dir,
			Overlay:    cfg.Overlay,
			Context:    cfg.Context,
			Logf:       cfg.Logf,
			Env:        cfg.Env,
			BuildFlags: cfg.BuildFlags,
			Tests:      cfg.Tests,
		},
	}
	return l
}

func (l *loader) Load(args ...string) (infov1.Project, error) {
	pkgs, err := packages.Load(l.cfg, args...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, ErrBadLoadRequest
	}
	if len(pkgs) == 0 {
		return nil, ErrPkgsEmpty
	}
	proj := make(map[*packages.Package]*infov1.Information, len(pkgs))
	for _, pkg := range pkgs {
		info, err := extractInfos(pkg, l.mngr.ParseMarkers)
		if err != nil {
			return nil, err
		}
		proj[pkg] = info
	}
	return proj, nil
}

func objectOf(pkg *packages.Package, ident *ast.Ident) (types.Object, error) {
	obj := pkg.TypesInfo.ObjectOf(ident)
	if obj == nil {
		return nil, fmt.Errorf("object not found: %v", ident)
	}
	return obj, nil
}

func typeOf(pkg *packages.Package, expr ast.Expr) (types.Type, error) {
	typ := pkg.TypesInfo.TypeOf(expr)
	if typ == nil {
		return nil, fmt.Errorf("type not found: %v", expr)
	}
	return typ, nil
}

func extractInfos(pkg *packages.Package, parse parseMarkers) (*infov1.Information, error) {
	info := newInformation()
	methodDecls := make([]*ast.FuncDecl, 0)
	for _, file := range pkg.Syntax {
		if err := extractFileInfo(pkg, parse, file, info); err != nil {
			return nil, err
		}
		for _, decl := range file.Decls {
			// this part of the code is only responsible for finding the correct
			// extract function.
			var err error
			var isTypeToken bool
			funcDecl, isFunc := decl.(*ast.FuncDecl)
			if isFunc && isMethod(funcDecl) {
				methodDecls = append(methodDecls, funcDecl)
				continue
			}
			if isFunc && !isMethod(funcDecl) {
				err = extractFuncInfo(pkg, parse, funcDecl, info)
			}
			if err != nil {
				return nil, err
			}
			if isFunc {
				continue
			}
			genDecl := decl.(*ast.GenDecl)
			switch genDecl.Tok {
			case token.CONST:
				err = extractConstInfo(pkg, parse, genDecl, info)
			case token.VAR:
				err = extractVarInfo(pkg, parse, genDecl, info)
			case token.IMPORT:
				err = extractImportInfo(pkg, parse, genDecl, info)
			default:
				isTypeToken = true
			}
			if err != nil {
				return nil, err
			}
			if !isTypeToken {
				continue
			}
			// TYPE token
			specs := convertSpecs[*ast.TypeSpec](genDecl.Specs)
			for _, spec := range specs {
				var err error
				typ, err := typeOf(pkg, spec.Name)
				if err != nil {
					return nil, err
				}
				// in go Aliases are not type Named
				_, isAlias := typ.(*types.Alias)
				if isAlias {
					err = extractAliasInfo(pkg, parse, genDecl, spec, info)
				}
				if err != nil {
					return nil, err
				}
				if isAlias {
					continue
				}
				switch typ.(*types.Named).Underlying().(type) {
				case *types.Struct:
					err = extractStructInfo(pkg, parse, genDecl, spec, info)
				case *types.Interface:
					err = extractIfaceInfo(pkg, parse, genDecl, spec, info)
				default:
					err = extractNamedInfo(pkg, parse, genDecl, spec, info)
				}
				if err != nil {
					return nil, err
				}
			}
		}
	}
	for _, methodDecl := range methodDecls {
		if err := extractMethodInfo(pkg, parse, methodDecl, info); err != nil {
			return nil, err
		}
	}

	return info, nil
}

func extractFileInfo(pkg *packages.Package, parse parseMarkers, file *ast.File, infos *infov1.Information) error {
	doc := file.Doc.Text()
	opts, err := parse(doc, optionv1.TargetPkg)
	if err != nil {
		return err
	}
	info := infov1.FileInfo{
		File: file,
		Opts: opts,
	}
	filename := filepath.Base(pkg.Fset.Position(file.Package).Filename)
	infos.Files[filename] = &info
	return nil
}

func extractMethodInfo(pkg *packages.Package, parse parseMarkers, decl *ast.FuncDecl, infos *infov1.Information) error {
	doc := decl.Doc.Text()
	opts, err := parse(doc, optionv1.TargetMethod)
	if err != nil {
		return err
	}
	info := infov1.FuncInfo{
		Decl: decl,
		Opts: opts,
	}
	rec := decl.Recv.List[0].Type
	recObj, err := objectOf(pkg, ident(rec))
	if err != nil {
		return err
	}
	methObj, err := objectOf(pkg, decl.Name)
	if err != nil {
		return err
	}
	return addMethodToType(recObj, methObj, info, infos)
}

func addMethodToType(receiver, method types.Object, methodInfo infov1.FuncInfo, info *infov1.Information) error {
	named, isNamed := info.Named[receiver]
	if isNamed {
		named.Methods[method] = &methodInfo
		return nil
	}
	strct, isStruct := info.Structs[receiver]
	if !isStruct {
		return fmt.Errorf("type is not extracted yet: %v", receiver)
	}
	strct.Methods[method] = &methodInfo
	return nil
}

func extractFuncInfo(pkg *packages.Package, parse parseMarkers, decl *ast.FuncDecl, infos *infov1.Information) error {
	doc := decl.Doc.Text()
	opts, err := parse(doc, optionv1.TargetFunc)
	if err != nil {
		return err
	}
	obj, err := objectOf(pkg, decl.Name)
	if err != nil {
		return err
	}
	info := infov1.FuncInfo{
		Decl: decl,
		Opts: opts,
	}
	infos.Funcs[obj] = &info
	return nil
}

func extractVarInfo(pkg *packages.Package, parse parseMarkers, decl *ast.GenDecl, infos *infov1.Information) error {
	specs := convertSpecs[*ast.ValueSpec](decl.Specs)
	for _, spec := range specs {
		doc := decl.Doc.Text() + spec.Doc.Text()
		for _, name := range spec.Names {
			opts, err := parse(doc, optionv1.TargetVar)
			if err != nil {
				return err
			}
			info := infov1.VarInfo{
				Spec: spec,
				Decl: decl,
				Opts: opts,
			}
			obj, err := objectOf(pkg, name)
			if err != nil {
				return err
			}
			infos.Vars[obj] = &info
		}
	}
	return nil
}

func extractConstInfo(pkg *packages.Package, parse parseMarkers, decl *ast.GenDecl, infos *infov1.Information) error {
	specs := convertSpecs[*ast.ValueSpec](decl.Specs)
	for _, spec := range specs {
		doc := decl.Doc.Text() + spec.Doc.Text()
		for _, name := range spec.Names {
			opts, err := parse(doc, optionv1.TargetConst)
			if err != nil {
				return err
			}
			info := infov1.ConstInfo{
				Spec: spec,
				Decl: decl,
				Opts: opts,
			}
			obj, err := objectOf(pkg, name)
			if err != nil {
				return err
			}
			infos.Consts[obj] = &info
		}
	}
	return nil
}

func extractImportInfo(pkg *packages.Package, parse parseMarkers, decl *ast.GenDecl, infos *infov1.Information) error {
	specs := convertSpecs[*ast.ImportSpec](decl.Specs)
	for _, spec := range specs {
		doc := decl.Doc.Text() + spec.Doc.Text()
		opts, err := parse(doc, optionv1.TargetImport)
		if err != nil {
			return err
		}
		info := infov1.ImportInfo{
			Spec: spec,
			Decl: decl,
			Opts: opts,
		}
		obj, err := objectOf(pkg, spec.Name)
		if err != nil {
			obj = pkg.TypesInfo.Implicits[spec]
		}
		if obj == nil {
			return fmt.Errorf("no types.Object found: %v", spec.Path.Value)
		}
		infos.Imports[obj] = &info
	}
	return nil
}

func extractAliasInfo(pkg *packages.Package, parse parseMarkers, decl *ast.GenDecl, spec *ast.TypeSpec, infos *infov1.Information) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	opts, err := parse(doc, optionv1.TargetAlias)
	if err != nil {
		return err
	}
	info := infov1.AliasInfo{
		Decl: decl,
		Spec: spec,
		Opts: opts,
	}
	obj, err := objectOf(pkg, spec.Name)
	if err != nil {
		return err
	}
	infos.Aliases[obj] = &info
	return nil
}

func extractNamedInfo(pkg *packages.Package, parse parseMarkers, decl *ast.GenDecl, spec *ast.TypeSpec, infos *infov1.Information) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	opts, err := parse(doc, optionv1.TargetNamed)
	if err != nil {
		return err
	}
	info := infov1.NamedInfo{
		Spec:    spec,
		Decl:    decl,
		Opts:    opts,
		Methods: make(map[types.Object]*infov1.FuncInfo),
	}
	obj, err := objectOf(pkg, spec.Name)
	if err != nil {
		return err
	}
	infos.Named[obj] = &info
	return nil
}

func extractStructInfo(pkg *packages.Package, parse parseMarkers, decl *ast.GenDecl, spec *ast.TypeSpec, infos *infov1.Information) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	opts, err := parse(doc, optionv1.TargetStruct)
	if err != nil {
		return err
	}
	obj, err := objectOf(pkg, spec.Name)
	if err != nil {
		return err
	}
	structType := spec.Type.(*ast.StructType)
	info := infov1.StructInfo{
		Opts:    opts,
		Spec:    spec,
		Decl:    decl,
		Fields:  make(map[types.Object]*infov1.FieldInfo, structType.Fields.NumFields()),
		Methods: make(map[types.Object]*infov1.FuncInfo, 0),
	}
	fieldInfos, err := fieldInfoOf(pkg, parse, structType)
	if err != nil {
		return err
	}
	info.Fields = fieldInfos
	infos.Structs[obj] = &info
	return nil
}

func fieldInfoOf(pkg *packages.Package, parse parseMarkers, spec *ast.StructType) (map[types.Object]*infov1.FieldInfo, error) {
	fields := make(map[types.Object]*infov1.FieldInfo, 0)
	for _, field := range spec.Fields.List {
		// embedded fields will be skipped
		if isEmbedded(field) {
			continue
		}
		doc := field.Doc.Text()
		opts, err := parse(doc, optionv1.TargetField)
		if err != nil {
			return nil, err
		}
		for _, name := range field.Names {
			info := infov1.FieldInfo{
				Ident: name,
				Field: field,
				Opts:  opts,
			}
			obj, err := objectOf(pkg, name)
			if err != nil {
				return nil, err
			}
			fields[obj] = &info
		}
	}
	return fields, nil
}

func extractIfaceInfo(pkg *packages.Package, parse parseMarkers, decl *ast.GenDecl, spec *ast.TypeSpec, infos *infov1.Information) error {
	doc := spec.Doc.Text() + decl.Doc.Text()
	opts, err := parse(doc, optionv1.TargetIface)
	if err != nil {
		return err
	}
	ifaceType := spec.Type.(*ast.InterfaceType)
	sigs, err := signatureInfoOf(pkg, parse, ifaceType)
	if err != nil {
		return err
	}
	info := infov1.IfaceInfo{
		Spec:       spec,
		Decl:       decl,
		Opts:       opts,
		Signatures: sigs,
	}
	obj, err := objectOf(pkg, spec.Name)
	if err != nil {
		return err
	}
	infos.Ifaces[obj] = &info
	return nil
}

func signatureInfoOf(pkg *packages.Package, parse parseMarkers, spec *ast.InterfaceType) (map[types.Object]*infov1.SignatureInfo, error) {
	sigs := make(map[types.Object]*infov1.SignatureInfo, spec.Methods.NumFields())
	for _, meth := range spec.Methods.List {
		doc := meth.Doc.Text()
		opts, err := parse(doc, optionv1.TargetIfaceSig)
		if err != nil {
			return nil, err
		}
		for _, name := range meth.Names {
			obj, err := objectOf(pkg, meth.Names[0])
			if err != nil {
				return nil, err
			}
			info := infov1.SignatureInfo{
				Ident:  name,
				Method: meth,
				Opts:   opts,
			}
			sigs[obj] = &info
		}
	}
	return sigs, nil
}
