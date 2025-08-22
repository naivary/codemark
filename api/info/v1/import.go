package v1

import "go/ast"

type ImportInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ImportSpec
	Opts Options
}

func (i *ImportInfo) Options() Options {
	return i.Opts
}
