package v1

import (
	"go/ast"
	"go/types"
)

type NamedInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Opts Options

	Methods map[types.Object]*FuncInfo
}

func (n *NamedInfo) Options() Options {
	return n.Opts
}
