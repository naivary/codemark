package v1

import "go/ast"

type FuncInfo struct {
	Decl *ast.FuncDecl
	Opts map[string][]any
}

func (f *FuncInfo) Options() map[string][]any {
	return f.Opts
}
