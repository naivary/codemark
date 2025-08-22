package v1

import "go/ast"

type VarInfo struct {
	Decl *ast.GenDecl
	Spec *ast.ValueSpec
	Opts Options
}

func (v *VarInfo) Options() Options {
	return v.Opts
}
