package v1

import "go/ast"

type ConstInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ValueSpec
	Opts map[string][]any
}

func (c *ConstInfo) Options() map[string][]any {
	return c.Opts
}
