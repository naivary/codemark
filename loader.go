package main

import (
	"fmt"
	"go/ast"
	"go/doc"
	gparser "go/parser"
	"go/token"
)

// Loader is responsible for loading the specified
// files and their documentation
type Loader interface {
	Load(files ...string) (*doc.Package, error)
}

func NewLoader() Loader {
	return &loader{}
}

var _ Loader = (*loader)(nil)

type loader struct{}

func (l *loader) Load(paths ...string) (*doc.Package, error) {
	fset := token.NewFileSet()
	files := make([]*ast.File, 0, len(paths))
	for _, path := range paths {
		file, err := gparser.ParseFile(fset, path, nil, gparser.ParseComments)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
        for _, decl := range file.Decls {
            fmt.Println(decl)
        }
	}
	return doc.NewFromFiles(fset, files, "codemark.com/loader")
}
