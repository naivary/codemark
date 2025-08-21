package v1

import "go/ast"

type VarInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ValueSpec
	Opts map[string][]any
}

func (v *VarInfo) Options() map[string][]any {
	return v.Opts
}
