package v1

import (
	"go/ast"
	"go/types"
)

type IfaceInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Opts Options

	Signatures map[types.Object]*SignatureInfo
}

func (i *IfaceInfo) Options() Options {
	return i.Opts
}

type SignatureInfo struct {
	Method *ast.Field
	Ident  *ast.Ident
	Opts   Options
}

func (s *SignatureInfo) Options() Options {
	return s.Opts
}
