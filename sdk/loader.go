package sdk

import (
	"errors"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

var (
	ErrPkgsEmpty = errors.New("loaded packages are empty. check that the defined patterns are matching any files")
)

type Project struct {
	Structs map[types.Object]*StructInfo
	Ifaces  map[types.Object]IfaceInfo
	Aliases map[types.Object]AliasInfo
	Named   map[types.Object]*NamedInfo
	Consts  map[types.Object]ConstInfo
	Vars    map[types.Object]VarInfo
	Imports map[types.Object]ImportInfo
	Funcs   map[types.Object]FuncInfo
	Pkgs    map[*packages.Package]PkgInfo
}

func NewProject() *Project {
	return &Project{
		Structs: make(map[types.Object]*StructInfo),
		Ifaces:  make(map[types.Object]IfaceInfo),
		Aliases: make(map[types.Object]AliasInfo),
		Named:   make(map[types.Object]*NamedInfo),
		Consts:  make(map[types.Object]ConstInfo),
		Vars:    make(map[types.Object]VarInfo),
		Imports: make(map[types.Object]ImportInfo),
		Funcs:   make(map[types.Object]FuncInfo),
		Pkgs:    make(map[*packages.Package]PkgInfo),
	}
}

// TODO: I think its not needed to make loader extensible because the existence
// of a project before loading it is required and must be provided independent
// of the loader.
// TODO: change return to map[*packages.Package]*Project
// and only ionclude one PkgInfo in the Prject return which includes all the
// markers of the packaghes over all the files
type Loader interface {
	Load(patterns ...string) ([]*Project, error)
}

type FuncInfo struct {
	Decl *ast.FuncDecl
	Pkg  *packages.Package
	Defs map[string][]any
}

type PkgInfo struct {
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

	Fields  map[types.Object]FieldInfo
	Methods map[types.Object]FuncInfo
}

type FieldInfo struct {
	Field *ast.Field
	Idn   *ast.Ident
	Defs  map[string][]any
}
