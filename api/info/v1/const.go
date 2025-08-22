package v1

import "go/ast"

type ConstInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ValueSpec
	Opts Options
}

func (c *ConstInfo) Options() Options {
	return c.Opts
}
