package utils

import (
	"fmt"
	"go/ast"
)

const ExprNameUndefined = "EXPR_NAME_UNDEFINED"

func ConvertSpecs[T any](specs []ast.Spec) []T {
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

func IsEmbedded(field *ast.Field) bool {
	return len(field.Names) == 0
}

func IsMethod(fn *ast.FuncDecl) bool {
	return fn.Recv != nil
}

func ExprName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return ExprName(t.X)
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", ExprName(t.X), t.Sel.Name)
	default:
		return ExprNameUndefined
	}
}
