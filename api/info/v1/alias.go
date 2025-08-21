package v1

import "go/ast"

type AliasInfo struct {
	Decl *ast.GenDecl
	Spec *ast.TypeSpec
	Opts map[string][]any
}

func (a *AliasInfo) Options() map[string][]any {
	return a.Opts
}
