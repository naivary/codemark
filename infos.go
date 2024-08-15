package main

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Info struct {
	Funcs      []*FuncInfo
	Methods    map[string][]*MethodInfo
	Consts     []*ConstInfo
	Vars       []*VarInfo
	Structs    []*StructInfo
	BasicTypes []*BasicTypeInfo
	Interfaces []*InterfaceInfo
	Aliases    []*AliasInfo
}

func NewInfo() *Info {
	return &Info{
		Methods: make(map[string][]*MethodInfo),
	}
}

type ConstInfo struct {
	// Name of the constant
	Name string
	// Documentation without the marker
	Doc   string
	Value ast.Expr
	Obj   types.Object
}

func NewConstInfo(spec *ast.ValueSpec, pkg *packages.Package) []*ConstInfo {
	infos := make([]*ConstInfo, 0, len(spec.Names))
	for i, ident := range spec.Names {
		value := spec.Values[i]
		obj := pkg.TypesInfo.Defs[ident]
		info := &ConstInfo{
			Name:  ident.Name,
			Value: value,
			Obj:   obj,
		}
		infos = append(infos, info)
	}
	return infos
}

type VarInfo struct {
	Name  string
	Doc   string
	Obj   types.Object
	Value ast.Expr
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
		}
		infos = append(infos, info)
	}
	return infos
}

type FuncInfo struct {
	Name string
	Doc  string
	Decl *ast.FuncDecl
}

func NewFuncInfo(decl *ast.FuncDecl) *FuncInfo {
	return &FuncInfo{
		Name: decl.Name.Name,
		Decl: decl,
	}
}

type MethodInfo struct {
	Name            string
	Doc             string
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
	}
}

type AliasInfo struct {
	Name string
	Doc  string
	Type *types.Alias
	Rhs  types.Type
	Decl *ast.GenDecl
}

func NewAliasInfo(spec *ast.TypeSpec, alias *types.Alias, decl *ast.GenDecl) *AliasInfo {
	return &AliasInfo{
		Name: spec.Name.Name,
		Rhs:  alias.Rhs(),
		Type: alias,
		Decl: decl,
	}
}

type InterfaceInfo struct {
	Name          string
	Doc           string
	Methods       []*FuncInfo
	Decl          *ast.GenDecl
	Type          *types.Interface
	InterfaceType *ast.InterfaceType
}

func NewInterfaceInfo(typeSpec *ast.TypeSpec, iface *types.Interface, decl *ast.GenDecl) *InterfaceInfo {
	ifaceType := typeSpec.Type.(*ast.InterfaceType)
	info := &InterfaceInfo{
		Name:          typeSpec.Name.Name,
		Decl:          decl,
		Type:          iface,
		InterfaceType: ifaceType,
		Methods:       make([]*FuncInfo, 0, len(ifaceType.Methods.List)),
	}
	for _, meth := range ifaceType.Methods.List {
		funcDecl := &ast.FuncDecl{
			Doc:  meth.Doc,
			Type: meth.Type.(*ast.FuncType),
			Name: meth.Names[0],
			Body: &ast.BlockStmt{},
			Recv: nil,
		}
		funcInfo := &FuncInfo{
			Name: meth.Names[0].Name,
			Decl: funcDecl,
		}
		info.Methods = append(info.Methods, funcInfo)
	}
	return info
}

type BasicTypeInfo struct {
	Name    string
	Doc     string
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
	}
	return info
}

type StructInfo struct {
	// Name of the Type.
	Name string
	// Doc string of the type without the markers.
	Doc string
	// Fields of the Type if it is a struct. If it's
	// not a struct it will be nil.
	Fields []*FieldInfo

	Type *types.Struct

	Decl *ast.GenDecl

	Methods []*MethodInfo

	Raw *ast.StructType
}

func NewStructInfo(typeSpec *ast.TypeSpec, strc *types.Struct, decl *ast.GenDecl, pkg *packages.Package) *StructInfo {
	structType := typeSpec.Type.(*ast.StructType)
	info := &StructInfo{
		Name: typeSpec.Name.Name,
		Type: strc,
		Decl: decl,
		Raw:  structType,
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
	Doc string

	Raw *ast.Field

	Type types.Type

	Obj types.Object

	Var *types.Var

	StarExpr *ast.StarExpr
}

func (f FieldInfo) IsEmbedded() bool {
	return f.Var.Embedded()
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
		}
		infos = append(infos, info)
	}
	return infos
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
