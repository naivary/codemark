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
	// indexed by filename
	Files map[string]FileInfo
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
		Files:   make(map[string]FileInfo),
	}
}

type Loader interface {
	Load(patterns ...string) (map[*packages.Package]*Project, error)
}

type FuncInfo struct {
	Decl *ast.FuncDecl
	Pkg  *packages.Package
	Defs map[string][]any
}

type FileInfo struct {
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
