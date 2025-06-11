package codemark

import (
	"fmt"
	"go/ast"
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

func exprName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return exprName(t.X)
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", t.X, t.Sel.Name)
	default:
		return ""
	}
}
