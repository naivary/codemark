package v1

import "go/ast"

type AliasInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Opts Options
}

func (a *AliasInfo) Options() Options {
	return a.Opts
}
