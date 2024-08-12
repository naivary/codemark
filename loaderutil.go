package main

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func isEmbedded(field *ast.Field) bool {
	return len(field.Names) == 0
}

func isMethod(list *ast.FieldList) bool {
	if list == nil {
		return false
	}
	return true
}

func isPosInFile(file *ast.File, pos token.Pos) bool {
	return file.FileStart <= pos && pos < file.FileEnd
}

func isVar(pkg *packages.Package, v *types.Var) bool {
	if v.IsField() {
		return false
	}

	// check if the variable is in a function/method declaration
	for _, file := range pkg.Syntax {
		pos := v.Pos()
		if !isPosInFile(file, pos) {
			// not in this file
			continue
		}
		path, _ := astutil.PathEnclosingInterval(file, pos, pos)
		for _, n := range path {
			switch n.(type) {
			case *ast.FuncDecl:
				return false
			}
		}
	}
	return true
}

// toNode is getting the node(s) of type `T` of the given position
func toNode[T any](files []*ast.File, pos token.Pos) []T {
	nodes := make([]T, 0, 0)
	for _, file := range files {
		if !isPosInFile(file, pos) {
			continue
		}
		path, _ := astutil.PathEnclosingInterval(file, pos, pos)
		for _, n := range path {
			v, ok := n.(T)
			if !ok {
				continue
			}
			nodes = append(nodes, v)
		}
	}
	return nodes
}
