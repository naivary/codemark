package main

import (
	"go/ast"
	"go/constant"
	"go/types"
)

type Info struct {
	Methods []*MethodInfo
	Funcs   []*FuncInfo
	Consts  []*ConstInfo
	Vars    []*VarInfo
	Types   []*TypeInfo
}

type MethodInfo struct {
	Name string
	Doc  string
}

type FuncInfo struct {
	Name string
	Doc  string
}

type ConstInfo struct {
	Name   string
	Doc    string
	Value  constant.Value
	Type   types.Type
	Object types.Object
	Ident  *ast.Ident
}

type VarInfo struct {
	Name string
	Doc  string
	Type types.Type
}

type TypeInfo struct {
	// Name of the Type.
	Name string
	// Doc string of the type without the markers.
	Doc string
	// IsStruct indicates if the type is a struct.
	IsStruct bool
	// Fields of the Type if it is a struct. If it's
	// not a struct it will be nil.
	Fields []*FieldInfo

	GenDecl *ast.GenDecl

	Type types.Type

	IsAlias bool

    IsBasic bool

    Ident *ast.Ident
}

func newTypeInfo(typeName *types.TypeName, decl *ast.GenDecl) *TypeInfo {
	return &TypeInfo{
		Name:    typeName.Name(),
		Type:    typeName.Type(),
		IsAlias: typeName.IsAlias(),
		GenDecl: decl,
	}
}

type FieldInfo struct {
	// Name of the field
	Name string
	// Doc string of the field
	Doc string

	IsEmbedded bool

	Type types.Type

	Expr ast.Expr

	Tags *ast.BasicLit

	Field *ast.Field
}

func newFieldInfo(field *ast.Field, typ types.Type) []*FieldInfo {
	infos := make([]*FieldInfo, 0, 0)
	if isEmbedded(field) {
		info := newEmbeddedField(field, typ)
		infos = append(infos, info)
		return infos
	}
	for _, idn := range field.Names {
		info := &FieldInfo{
			Name:  idn.Name,
			Doc:   field.Doc.Text(),
			Field: field,
			Tags:  field.Tag,
			Type:  typ,
			Expr:  field.Type,
		}
		infos = append(infos, info)
	}
	return infos
}

func newEmbeddedField(field *ast.Field, typ types.Type) *FieldInfo {
	name := field.Type.(*ast.Ident).Name
	return &FieldInfo{
		Name:       name,
		Doc:        field.Doc.Text(),
		IsEmbedded: true,
		Field:      field,
		Tags:       field.Tag,
		Expr:       field.Type,
		Type:       typ,
	}
}
