package v1

import (
	"go/ast"
	"go/types"
)

type IfaceInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Opts map[string][]any

	Signatures map[types.Object]*SignatureInfo
}

func (i *IfaceInfo) Options() map[string][]any {
	return i.Opts
}

type SignatureInfo struct {
	Method *ast.Field
	Ident  *ast.Ident
	Opts   map[string][]any
}

func (s *SignatureInfo) Options() map[string][]any {
	return s.Opts
}
