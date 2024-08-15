package main

import "go/ast"

func convertSpecs[T any](specs []ast.Spec) []T {
	converted := make([]T, 0, len(specs))
	for _, spec := range specs {
		v, ok := spec.(T)
		if !ok {
			return nil
		}
		converted = append(converted, v)
	}
	return converted
}

func isEmbedded(field *ast.Field) bool {
	return len(field.Names) == 0
}
