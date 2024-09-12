package main

import (
	"go/ast"

	"github.com/naivary/codemark/parser"
)

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

func isMethod(fn *ast.FuncDecl) bool {
	return fn.Recv != nil
}

func newDefinitions(doc string, t Target, conv Converter) (Definitions, error) {
	markers, err := parser.Parse(doc)
	if err != nil {
		return nil, err
	}
	defs := Definitions{}
	for _, marker := range markers {
		def, err := conv.Convert(marker, t)
		if err != nil {
			return nil, err
		}
		idn := marker.Ident()
		defs.add(idn, def)
	}
	return defs, nil
}
