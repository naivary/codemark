package v1

import "go/ast"

type ImportInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ImportSpec
	Opts map[string][]any
}

func (i *ImportInfo) Options() map[string][]any {
	return i.Opts
}
