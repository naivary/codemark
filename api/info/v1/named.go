package v1

import (
	"go/ast"
	"go/types"
)

type NamedInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Opts map[string][]any

	Methods map[types.Object]*FuncInfo
}

func (n *NamedInfo) Options() map[string][]any {
	return n.Opts
}
