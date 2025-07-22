package laoder

import (
	"go/ast"
	"go/types"
)

type Optioner interface {
	Options() map[string][]any
}

type Filename = string

type Information struct {
	Structs map[types.Object]*StructInfo
	Ifaces  map[types.Object]IfaceInfo
	Aliases map[types.Object]AliasInfo
	Named   map[types.Object]*NamedInfo
	Consts  map[types.Object]ConstInfo
	Vars    map[types.Object]VarInfo
	Imports map[types.Object]ImportInfo
	Funcs   map[types.Object]FuncInfo
	Files   map[Filename]FileInfo
}

type FuncInfo struct {
	Decl *ast.FuncDecl
	Opts map[string][]any
}

func (f FuncInfo) Options() map[string][]any {
	return f.Opts
}

type FileInfo struct {
	File *ast.File
	Opts map[string][]any
}

func (f FileInfo) Options() map[string][]any {
	return f.Opts
}

type ImportInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ImportSpec
	Opts map[string][]any
}

func (i ImportInfo) Options() map[string][]any {
	return i.Opts
}

type VarInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ValueSpec
	Opts map[string][]any
}

func (v VarInfo) Options() map[string][]any {
	return v.Opts
}

type ConstInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ValueSpec
	Opts map[string][]any
}

func (c ConstInfo) Options() map[string][]any {
	return c.Opts
}

type NamedInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Opts map[string][]any

	Methods map[types.Object]FuncInfo
}

func (n NamedInfo) Options() map[string][]any {
	return n.Opts
}

type AliasInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Opts map[string][]any
}

func (a AliasInfo) Options() map[string][]any {
	return a.Opts
}

type IfaceInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Opts map[string][]any

	Signatures map[types.Object]SignatureInfo
}

func (i IfaceInfo) Options() map[string][]any {
	return i.Opts
}

type SignatureInfo struct {
	Method *ast.Field
	Ident  *ast.Ident
	Opts   map[string][]any
}

func (s SignatureInfo) Options() map[string][]any {
	return s.Opts
}

type StructInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Opts map[string][]any

	Fields  map[types.Object]FieldInfo
	Methods map[types.Object]FuncInfo
}

func (s StructInfo) Options() map[string][]any {
	return s.Opts
}

type FieldInfo struct {
	Field *ast.Field
	Ident *ast.Ident
	Opts  map[string][]any
}

func (f FieldInfo) Options() map[string][]any {
	return f.Opts
}
