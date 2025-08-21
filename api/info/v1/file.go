package v1

import "go/ast"

type FileInfo struct {
	File *ast.File
	Opts map[string][]any
}

func (f *FileInfo) Options() map[string][]any {
	return f.Opts
}
