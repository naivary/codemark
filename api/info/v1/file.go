package v1

import "go/ast"

type FileInfo struct {
	File *ast.File
	Opts Options
}

func (f *FileInfo) Options() Options {
	return f.Opts
}
