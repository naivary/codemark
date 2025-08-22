package v1

import "go/ast"

type FuncInfo struct {
	Decl *ast.FuncDecl
	Opts Options
}

func (f *FuncInfo) Options() Options {
	return f.Opts
}
