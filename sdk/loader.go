package sdk

import (
	"errors"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

var (
	ErrPkgsEmpty = errors.New("loaded packages are empty. check that the defined patterns are matching any files")
)

type Project struct {
	Structs []*StructInfo
	Ifaces  []IfaceInfo
	Aliases []AliasInfo
	Named   []*NamedInfo
	Consts  []ConstInfo
	Vars    []VarInfo
	Imports []ImportInfo
	Funcs   []FuncInfo
	Pkgs    []PkgInfo
}

type Loader interface {
	Load(patterns ...string) ([]*Project, error)
}

type FuncInfo struct {
	Decl *ast.FuncDecl
	Pkg  *packages.Package
	Defs map[string][]any
}

type PkgInfo struct {
	Pkg  *packages.Package
	File *ast.File
	Defs map[string][]any
}

type ImportInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ImportSpec
	Pkg  *packages.Package
	Defs map[string][]any
}

type VarInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ValueSpec
	Pkg  *packages.Package
	Defs map[string][]any
}

type ConstInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ValueSpec
	Pkg  *packages.Package
	Defs map[string][]any
}

type NamedInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Pkg  *packages.Package
	Defs map[string][]any

	Methods []FuncInfo
}

type AliasInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Pkg  *packages.Package
	Defs map[string][]any
}

type IfaceInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Pkg  *packages.Package
	Defs map[string][]any

	Signatures []SignatureInfo
}

type SignatureInfo struct {
	Method *ast.Field
	Idn    *ast.Ident
	Defs   map[string][]any
}

type StructInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Pkg  *packages.Package
	Defs map[string][]any

	Fields  []FieldInfo
	Methods []FuncInfo
}

type FieldInfo struct {
	Field *ast.Field
	Idn   *ast.Ident
	Defs  map[string][]any
}
