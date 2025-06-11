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
	Structs []StructInfo
}

type Loader interface {
	Load(patterns ...string) (*Project, error)
}

type FieldInfo struct {
	Field *ast.Field
	Idn   *ast.Ident
	Defs  map[string][]any
}

type StructInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Pkg  *packages.Package
	Defs map[string][]any

	Fields []FieldInfo
}

func (s *StructInfo) String() string {
	return s.Name()
}

func (s *StructInfo) Name() string {
	// return the name
	return ""
}
