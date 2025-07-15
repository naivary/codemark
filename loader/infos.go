package loader

import (
	"errors"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

var (
	ErrPkgsEmpty = errors.New("loaded packages are empty. check that the defined patterns are matching any files")
)

type Filename = string

// TODO: Find a better name than Project
type Project struct {
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
		Files:   make(map[Filename]FileInfo),
	}
}

type Loader interface {
	Load(patterns ...string) (map[*packages.Package]*Project, error)
}

type FuncInfo struct {
	Decl *ast.FuncDecl
	Defs map[string][]any
}

func (f FuncInfo) Definitions() map[string][]any {
	return f.Defs
}

type FileInfo struct {
	File *ast.File
	Defs map[string][]any
}

func (f FileInfo) Definitions() map[string][]any {
	return f.Defs
}

type ImportInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ImportSpec
	Defs map[string][]any
}

func (i ImportInfo) Definitions() map[string][]any {
	return i.Defs
}

type VarInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ValueSpec
	Defs map[string][]any
}

func (v VarInfo) Definitions() map[string][]any {
	return v.Defs
}

type ConstInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ValueSpec
	Defs map[string][]any
}

func (c ConstInfo) Definitions() map[string][]any {
	return c.Defs
}

type NamedInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Defs map[string][]any

	Methods map[types.Object]FuncInfo
}

func (n NamedInfo) Definitions() map[string][]any {
	return n.Defs
}

type AliasInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Defs map[string][]any
}

func (a AliasInfo) Definitions() map[string][]any {
	return a.Defs
}

type IfaceInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Defs map[string][]any

	Signatures map[types.Object]SignatureInfo
}

func (i IfaceInfo) Definitions() map[string][]any {
	return i.Defs
}

type SignatureInfo struct {
	Method *ast.Field
	Idn    *ast.Ident
	Defs   map[string][]any
}

func (s SignatureInfo) Definitions() map[string][]any {
	return s.Defs
}

type StructInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Defs map[string][]any

	Fields  map[types.Object]FieldInfo
	Methods map[types.Object]FuncInfo
}

func (s StructInfo) Definitions() map[string][]any {
	return s.Defs
}

type FieldInfo struct {
	Field *ast.Field
	Idn   *ast.Ident
	Defs  map[string][]any
}

func (f FieldInfo) Definitions() map[string][]any {
	return f.Defs
}
