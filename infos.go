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
	Fields []FieldInfo

	GenDecl *ast.GenDecl

	Type types.Type

    Markers 
}

type FieldInfo struct {
	Name string
	Doc  string
}
