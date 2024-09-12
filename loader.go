package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// TODO: I have to include eveything from types.* e.g. *types.Var or
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
			err := l.funcDecl(pkg, funcDecl)
			if err != nil {
				return err
			}
			continue
		}
		var err error
		genDecl := decl.(*ast.GenDecl)
		switch genDecl.Tok {
		case token.CONST:
			err = l.constInfo(pkg, genDecl)
		case token.VAR:
			err = l.varInfo(pkg, genDecl)
		case token.IMPORT:
			err = l.importInfo(pkg, genDecl)
		case token.TYPE:
			err = l.typeDecl(pkg, genDecl)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *loader) funcDecl(pkg *packages.Package, fn *ast.FuncDecl) error {
	if isMethod(fn) {
		return l.methodInfo(pkg, fn)
	}
	return l.funcInfo(pkg, fn)
}

func (l *loader) typeDecl(pkg *packages.Package, gen *ast.GenDecl) error {
	specs := convertSpecs[*ast.TypeSpec](gen.Specs)
	for _, spec := range specs {
		typ := pkg.TypesInfo.TypeOf(spec.Name)
		alias, isAlias := typ.(*types.Alias)
		if isAlias {
			err := l.aliasInfo(pkg, gen, alias, spec)
			if err != nil {
				return err
			}
			continue
		}
		named := typ.(*types.Named).Underlying()
		strct, isStruct := named.(*types.Struct)
		if isStruct {
			err := l.structInfo(pkg, gen, strct, spec)
			if err != nil {
				return err
			}
			continue
		}
		iface, isInterface := named.(*types.Interface)
		if isInterface {
			err := l.interfaceInfo(pkg, gen, iface, spec)
			if err != nil {
				return err
			}
			continue
		}
		err := l.typeInfo(pkg, gen, named, spec)
		if err != nil {
			return err
		}
		continue

	}
	return nil
}

func (l *loader) funcInfo(pkg *packages.Package, fn *ast.FuncDecl) error {
	typ := pkg.TypesInfo.TypeOf(fn.Type)
	obj := pkg.TypesInfo.ObjectOf(fn.Name)
	doc := fn.Doc.Text()
	defs, err := newDefinitions(doc, TargetFunc, l.conv)
	if err != nil {
		return err
	}
	info := &FuncInfo{
		doc:  doc,
		defs: defs,
		Decl: fn,
		typ:  typ,
		obj:  obj,
	}
	l.pkgInfo.Funcs = append(l.pkgInfo.Funcs, info)
	return nil
}

func (l *loader) methodInfo(pkg *packages.Package, meth *ast.FuncDecl) error {
	typ := pkg.TypesInfo.TypeOf(meth.Type)
	obj := pkg.TypesInfo.ObjectOf(meth.Name)
	doc := meth.Doc.Text()
	defs, err := newDefinitions(doc, TargetMethod, l.conv)
	if err != nil {
		return err
	}
	info := &MethodInfo{
		doc:  doc,
		defs: defs,
		obj:  obj,
		typ:  typ,
		Decl: meth,
	}
	l.pkgInfo.Methods = append(l.pkgInfo.Methods, info)
	return nil
}

func (l *loader) constInfo(pkg *packages.Package, gen *ast.GenDecl) error {
	specs := convertSpecs[*ast.ValueSpec](gen.Specs)
	for _, spec := range specs {
		doc := gen.Doc.Text() + spec.Doc.Text()
		defs, err := newDefinitions(doc, TargetConst, l.conv)
		if err != nil {
			return err
		}
		for i, idn := range spec.Names {
			var value ast.Expr
			if len(spec.Values) > 0 {
				value = spec.Values[i]
			}
			typ := pkg.TypesInfo.TypeOf(value)
			obj := pkg.TypesInfo.ObjectOf(idn)
			info := &ConstInfo{
				doc:   doc,
				idn:   idn,
				value: value,
				typ:   typ,
				obj:   obj,
				defs:  defs,
			}
			l.pkgInfo.Consts = append(l.pkgInfo.Consts, info)
		}
	}
	return nil
}

func (l *loader) varInfo(pkg *packages.Package, gen *ast.GenDecl) error {
	specs := convertSpecs[*ast.ValueSpec](gen.Specs)
	for _, spec := range specs {
		doc := gen.Doc.Text() + spec.Doc.Text()
		defs, err := newDefinitions(doc, TargetVar, l.conv)
		if err != nil {
			return err
		}
		for i, idn := range spec.Names {
			var value ast.Expr
			if len(spec.Values) > 0 {
				value = spec.Values[i]
			}
			typ := pkg.TypesInfo.TypeOf(value)
			obj := pkg.TypesInfo.ObjectOf(idn)
			info := &VarInfo{
				doc:   doc,
				idn:   idn,
				value: value,
				typ:   typ,
				obj:   obj,
				defs:  defs,
			}
			l.pkgInfo.Vars = append(l.pkgInfo.Vars, info)
		}
	}
	return nil
}

func (l *loader) structInfo(pkg *packages.Package, gen *ast.GenDecl, strct *types.Struct, spec *ast.TypeSpec) error {
	structType := spec.Type.(*ast.StructType)
	doc := spec.Doc.Text() + gen.Doc.Text()
	defs, err := newDefinitions(doc, TargetType, l.conv)
	if err != nil {
		return err
	}
	info := &StructInfo{
		idn:    spec.Name,
		doc:    doc,
		defs:   defs,
		typ:    strct,
		obj:    pkg.TypesInfo.ObjectOf(spec.Name),
		spec:   spec,
		fields: make([]*FieldInfo, 0, len(structType.Fields.List)),
	}

	for _, field := range structType.Fields.List {
		if isEmbedded(field) {
			typ := pkg.TypesInfo.TypeOf(field.Type)
			doc := field.Doc.Text()
			defs, err := newDefinitions(doc, TargetField, l.conv)
			if err != nil {
				return err
			}
			fieldInfo := &FieldInfo{
				doc:  doc,
				expr: field.Type,
				typ:  typ,
				defs: defs,
			}
			info.fields = append(info.fields, fieldInfo)
			continue
		}

		for _, idn := range field.Names {
			typ := pkg.TypesInfo.TypeOf(field.Type)
			obj := pkg.TypesInfo.ObjectOf(idn)
			doc := field.Doc.Text()
			defs, err := newDefinitions(doc, TargetField, l.conv)
			if err != nil {
				return err
			}
			fieldInfo := &FieldInfo{
				doc:  doc,
				idn:  idn,
				expr: field.Type,
				typ:  typ,
				obj:  obj,
				defs: defs,
			}
			info.fields = append(info.fields, fieldInfo)
		}
	}
	l.pkgInfo.Structs = append(l.pkgInfo.Structs, info)
	return nil
}

func (l *loader) interfaceInfo(pkg *packages.Package, gen *ast.GenDecl, iface *types.Interface, spec *ast.TypeSpec) error {
	ifaceType := spec.Type.(*ast.InterfaceType)
	typ := pkg.TypesInfo.TypeOf(spec.Type)
	obj := pkg.TypesInfo.ObjectOf(spec.Name)
	doc := gen.Doc.Text() + spec.Doc.Text()
	defs, err := newDefinitions(doc, TargetInterface, l.conv)
	if err != nil {
		return err
	}
	info := &InterfaceInfo{
		doc:  doc,
		idn:  spec.Name,
		typ:  typ,
		obj:  obj,
		defs: defs,
	}

	for _, meth := range ifaceType.Methods.List {
		if isEmbedded(meth) {
			// TODO: for now embedded fields will be ignored but there should be
			// a consistent rule how to integreate them as signature
			// informations. Maybe a slice which is containing the other
			// interfaces found?
			continue
		}
		doc := meth.Doc.Text()
		defs, err := newDefinitions(doc, TargetInterfaceSignature, l.conv)
		if err != nil {
			return err
		}
		name := meth.Names[0]
		obj := pkg.TypesInfo.ObjectOf(name)
		typ := pkg.TypesInfo.TypeOf(meth.Type)
		signatureInfo := &SignatureInfo{
			doc:  doc,
			idn:  name,
			typ:  typ,
			obj:  obj,
			defs: defs,
		}
		info.signatures = append(info.signatures, signatureInfo)
	}
	l.pkgInfo.Interfaces = append(l.pkgInfo.Interfaces, info)
	return nil

}

func (l *loader) aliasInfo(pkg *packages.Package, gen *ast.GenDecl, alias *types.Alias, spec *ast.TypeSpec) error {
	typ := pkg.TypesInfo.TypeOf(spec.Type)
	obj := pkg.TypesInfo.ObjectOf(spec.Name)
	doc := gen.Doc.Text() + spec.Doc.Text()
	defs, err := newDefinitions(doc, TargetAlias, l.conv)
	if err != nil {
		return err
	}
	info := &AliasInfo{
		doc:   doc,
		typ:   typ,
		obj:   obj,
		idn:   spec.Name,
		defs:  defs,
		alias: alias,
	}
	l.pkgInfo.Aliases = append(l.pkgInfo.Aliases, info)
	return nil
}

func (l *loader) importInfo(pkg *packages.Package, gen *ast.GenDecl) error {
	specs := convertSpecs[*ast.ImportSpec](gen.Specs)
	importDoc := gen.Doc.Text()
	defs, err := newDefinitions(importDoc, TargetImport, l.conv)
	if err != nil {
		return err
	}
	info := &ImportInfo{
		defs: defs,
		doc:  importDoc,
	}

	for _, spec := range specs {
		doc := spec.Doc.Text()
		defs, err := newDefinitions(doc, TargetImportedPackage, l.conv)
		if err != nil {
			return err
		}
		importedPkgInfo := &ImportedPackageInfo{
			doc:  doc,
			defs: defs,
			spec: spec,
		}
		info.pkgs = append(info.pkgs, importedPkgInfo)
	}
	l.pkgInfo.Imports = append(l.pkgInfo.Imports, info)
	return nil
}

func (l *loader) typeInfo(pkg *packages.Package, gen *ast.GenDecl, typ types.Type, spec *ast.TypeSpec) error {
	pkgTyp := pkg.TypesInfo.TypeOf(spec.Type)
	obj := pkg.TypesInfo.ObjectOf(spec.Name)
	doc := gen.Doc.Text() + spec.Doc.Text()
	defs, err := newDefinitions(doc, TargetType, l.conv)
	if err != nil {
		return err
	}
	info := &TypeInfo{
		idn:  spec.Name,
		typ:  pkgTyp,
		obj:  obj,
		doc:  doc,
		defs: defs,
	}
	l.pkgInfo.Types = append(l.pkgInfo.Types, info)
	return nil
}
