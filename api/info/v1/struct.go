package v1

import (
	"go/ast"
	"go/types"
)

type StructInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Opts Options

	Fields  map[types.Object]*FieldInfo
	Methods map[types.Object]*FuncInfo
}

func (s *StructInfo) Options() Options {
	return s.Opts
}

func (s *StructInfo) HasField(ident string) bool {
	for _, field := range s.Fields {
		if field.Ident.Name == ident {
			return true
		}
	}
	return false
}

func (s *StructInfo) GetField(ident string) *FieldInfo {
	for _, field := range s.Fields {
		if field.Ident.Name == ident {
			return field
		}
	}
	return nil
}

type FieldInfo struct {
	Field *ast.Field
	Ident *ast.Ident
	Opts  Options
}

func (f *FieldInfo) Options() Options {
	return f.Opts
}
