package main

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type DocDefer interface {
	Doc() string
	Defs() Definitions
}

type Definitions map[string][]any

func (d Definitions) Add(idn string, def any) {
	defs, ok := d[idn]
	if !ok {
		d[idn] = []any{def}
		return
	}
	d[idn] = append(defs, def)
}

func (d Definitions) IsSet(idn string) bool {
	_, ok := d[idn]
	return ok
}

type Info struct {
	Consts     []*ConstInfo
	Vars       []*VarInfo
	Funcs      []*FuncInfo
	Methods    map[string][]*MethodInfo
	BasicTypes []*BasicTypeInfo
	Structs    []*StructInfo
	Interfaces []*InterfaceInfo
	Aliases    []*AliasInfo
	Imports    []*ImportStmtInfo
}

func NewInfo() *Info {
	return &Info{
		Methods: make(map[string][]*MethodInfo),
	}
}

type ImportInfo struct {
	Name string
	doc  string
	Spec *ast.ImportSpec
	defs Definitions
}

func NewImportInfo(spec *ast.ImportSpec) *ImportInfo {
	name := spec.Name.String()
	if name == "" {
		name = spec.Path.Value
	}
	return &ImportInfo{
		Name: name,
		Spec: spec,
	}
}

func (i ImportInfo) Doc() string {
	return i.doc
}

func (i ImportInfo) Defs() Definitions {
	return i.defs
}

type ImportStmtInfo struct {
	Name    string
	doc     string
	Decl    *ast.GenDecl
	Imports []*ImportInfo
	defs    Definitions
}

func NewImportStmtInfo(decl *ast.GenDecl) *ImportStmtInfo {
	specs := convertSpecs[*ast.ImportSpec](decl.Specs)
	info := &ImportStmtInfo{}
	for _, spec := range specs {
		info.Imports = append(info.Imports, NewImportInfo(spec))
	}
	return info
}

func (i ImportStmtInfo) Doc() string {
	return i.doc
}

func (i ImportStmtInfo) Defs() Definitions {
	return i.defs
}

type ConstInfo struct {
	Name  string
	doc   string
	Value ast.Expr
	Obj   types.Object
	defs  Definitions
}

func NewConstInfo(spec *ast.ValueSpec, pkg *packages.Package) []*ConstInfo {
	infos := make([]*ConstInfo, 0, len(spec.Names))
	for i, ident := range spec.Names {
		value := spec.Values[i]
		obj := pkg.TypesInfo.Defs[ident]
		info := &ConstInfo{
			Name:  ident.Name,
			doc:   spec.Doc.Text(),
			Value: value,
			Obj:   obj,
			defs:  Definitions{},
		}
		infos = append(infos, info)
	}
	return infos
}

func (c *ConstInfo) Doc() string {
	return c.doc
}

func (c *ConstInfo) Defs() Definitions {
	return c.defs
}

type VarInfo struct {
	Name  string
	doc   string
	Obj   types.Object
	Value ast.Expr

	defs Definitions
}

func NewVarInfo(spec *ast.ValueSpec, pkg *packages.Package) []*VarInfo {
	infos := make([]*VarInfo, 0, len(spec.Names))
	for i, ident := range spec.Names {
		value := spec.Values[i]
		obj := pkg.TypesInfo.Defs[ident]
		info := &VarInfo{
			Name:  ident.Name,
			Value: value,
			Obj:   obj,
			doc:   spec.Doc.Text(),
			defs:  Definitions{},
		}
		infos = append(infos, info)
	}
	return infos
}

func (v VarInfo) Doc() string {
	return v.doc
}

func (v VarInfo) Defs() Definitions {
	return v.defs
}

type FuncInfo struct {
	Name string
	doc  string
	defs Definitions
	Decl *ast.FuncDecl
}

func NewFuncInfo(decl *ast.FuncDecl) *FuncInfo {
	return &FuncInfo{
		Name: decl.Name.Name,
		Decl: decl,
		doc:  decl.Doc.Text(),
		defs: Definitions{},
	}
}

func (f FuncInfo) Doc() string {
	return f.doc
}

func (f FuncInfo) Defs() Definitions {
	return f.defs
}

type MethodInfo struct {
	Name            string
	doc             string
	defs            Definitions
	Decl            *ast.FuncDecl
	ValueReceiver   *ast.Ident
	PointerReceiver *ast.StarExpr
	Expr            ast.Expr
}

func NewMethodInfo(decl *ast.FuncDecl, valueRec *ast.Ident, ptrRec *ast.StarExpr) *MethodInfo {
	return &MethodInfo{
		Name:            decl.Name.Name,
		Decl:            decl,
		ValueReceiver:   valueRec,
		PointerReceiver: ptrRec,
		Expr:            decl.Recv.List[0].Type,
		defs:            Definitions{},
	}
}

func (m MethodInfo) Doc() string {
	return m.doc
}

func (m MethodInfo) Defs() Definitions {
	return m.defs
}

type AliasInfo struct {
	Name string
	doc  string
	Type *types.Alias
	Rhs  types.Type
	Decl *ast.GenDecl

	defs Definitions
}

func NewAliasInfo(spec *ast.TypeSpec, alias *types.Alias, decl *ast.GenDecl) *AliasInfo {
	return &AliasInfo{
		Name: spec.Name.Name,
		Rhs:  alias.Rhs(),
		Type: alias,
		Decl: decl,
		doc:  spec.Doc.Text(),
		defs: Definitions{},
	}
}

func (a AliasInfo) Doc() string {
	return a.doc
}

func (a AliasInfo) Defs() Definitions {
	return a.defs
}

type InterfaceInfo struct {
	Name          string
	doc           string
	defs          Definitions
	Funcs         []*InterfaceFuncInfo
	Decl          *ast.GenDecl
	Type          *types.Interface
	InterfaceType *ast.InterfaceType
}

type InterfaceFuncInfo struct {
	doc  string
	defs Definitions
	Name *ast.Ident
	Type *ast.FuncType
}

func (i InterfaceFuncInfo) Doc() string {
	return i.doc
}

func (i InterfaceFuncInfo) Defs() Definitions {
	return i.defs
}

func NewInterfaceInfo(typeSpec *ast.TypeSpec, iface *types.Interface, decl *ast.GenDecl) *InterfaceInfo {
	ifaceType := typeSpec.Type.(*ast.InterfaceType)
	info := &InterfaceInfo{
		Name:          typeSpec.Name.Name,
		Decl:          decl,
		Type:          iface,
		InterfaceType: ifaceType,
		Funcs:         make([]*InterfaceFuncInfo, 0, len(ifaceType.Methods.List)),
		defs:          Definitions{},
		doc:           typeSpec.Doc.Text(),
	}
	for _, meth := range ifaceType.Methods.List {
		ifaceFuncInfo := &InterfaceFuncInfo{
			Name: meth.Names[0],
			doc:  meth.Doc.Text(),
			Type: meth.Type.(*ast.FuncType),
			defs: Definitions{},
		}
		info.Funcs = append(info.Funcs, ifaceFuncInfo)
	}
	return info
}

func (i InterfaceInfo) Doc() string {
	return i.doc
}

func (i InterfaceInfo) Defs() Definitions {
	return i.defs
}

type BasicTypeInfo struct {
	Name    string
	doc     string
	defs    Definitions
	Type    *types.Basic
	Pointer *types.Pointer
	Decl    *ast.GenDecl
	Methods []*MethodInfo
}

func NewBasicTypeInfo(typeSpec *ast.TypeSpec, basic *types.Basic, decl *ast.GenDecl, ptr *types.Pointer) *BasicTypeInfo {
	info := &BasicTypeInfo{
		Name:    typeSpec.Name.Name,
		Type:    basic,
		Decl:    decl,
		Pointer: ptr,
		defs:    Definitions{},
		doc:     typeSpec.Doc.Text(),
	}
	return info
}

func (b BasicTypeInfo) Doc() string {
	return b.doc
}

func (b BasicTypeInfo) Defs() Definitions {
	return b.defs
}

type StructInfo struct {
	// Name of the Type.
	Name string
	// Doc string of the type without the markers.
	doc string
	// Fields of the Type if it is a struct. If it's
	// not a struct it will be nil.
	Fields []*FieldInfo

	Type *types.Struct

	Decl *ast.GenDecl

	Methods []*MethodInfo

	Raw *ast.StructType

	defs Definitions
}

func (s StructInfo) Doc() string {
	return s.doc
}

func (s StructInfo) Defs() Definitions {
	return s.defs
}

func NewStructInfo(typeSpec *ast.TypeSpec, strc *types.Struct, decl *ast.GenDecl, pkg *packages.Package) *StructInfo {
	structType := typeSpec.Type.(*ast.StructType)
	info := &StructInfo{
		Name: typeSpec.Name.Name,
		Type: strc,
		Decl: decl,
		Raw:  structType,
		doc:  typeSpec.Doc.Text(),
		defs: Definitions{},
	}
	for _, field := range structType.Fields.List {
		fieldInfos := NewFieldInfo(field, pkg)
		info.Fields = append(info.Fields, fieldInfos...)
	}
	return info
}

type FieldInfo struct {
	// Name of the field
	Name string
	// Doc string of the field
	doc string

	defs Definitions

	Raw *ast.Field

	Type types.Type

	Obj types.Object

	Var *types.Var

	StarExpr *ast.StarExpr
}

func NewFieldInfo(field *ast.Field, pkg *packages.Package) []*FieldInfo {
	infos := make([]*FieldInfo, 0, len(field.Names))
	if isEmbedded(field) {
		infos = append(infos, newEmbeddedFieldInfo(field, pkg))
		return infos
	}
	for _, ident := range field.Names {
		typ := pkg.TypesInfo.TypeOf(field.Type)
		obj := pkg.TypesInfo.ObjectOf(ident)
		info := &FieldInfo{
			Name: ident.Name,
			Type: typ,
			Raw:  field,
			Var:  obj.(*types.Var),
			Obj:  obj,
			doc:  field.Doc.Text(),
			defs: Definitions{},
		}
		infos = append(infos, info)
	}
	return infos
}

func (f FieldInfo) Doc() string {
	return f.doc
}

func (f FieldInfo) Defs() Definitions {
	return f.defs
}

func (f FieldInfo) IsEmbedded() bool {
	return f.Var.Embedded()
}

func newEmbeddedFieldInfo(field *ast.Field, pkg *packages.Package) *FieldInfo {
	var ptr *ast.StarExpr
	ident, isIdent := field.Type.(*ast.Ident)
	if !isIdent {
		ptr = field.Type.(*ast.StarExpr)
		ident = ptr.X.(*ast.Ident)
	}
	obj := pkg.TypesInfo.ObjectOf(ident)
	return &FieldInfo{
		Name:     ident.Name,
		Raw:      field,
		Type:     pkg.TypesInfo.TypeOf(field.Type),
		Obj:      obj,
		Var:      obj.(*types.Var),
		StarExpr: ptr,
	}
}
