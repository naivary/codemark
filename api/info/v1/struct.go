package v1

import (
	"go/ast"
	"go/types"
)

type StructInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Opts map[string][]any

	Fields  map[types.Object]*FieldInfo
	Methods map[types.Object]*FuncInfo
}

func (s *StructInfo) Options() map[string][]any {
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

func (s *StructInfo) GetFieldByIdent(ident string) *FieldInfo {
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
	Opts  map[string][]any
}

func (f *FieldInfo) Options() map[string][]any {
	return f.Opts
}
