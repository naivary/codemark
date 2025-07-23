package loader

import (
	"go/ast"
)

func _map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func convertSpecs[T any](specs []ast.Spec) []T {
	return _map(specs, func(spec ast.Spec) T {
		return spec.(T)
	})
}

func isEmbedded(field *ast.Field) bool {
	return len(field.Names) == 0
}

func isMethod(fn *ast.FuncDecl) bool {
	return fn.Recv != nil
}

func ident(expr ast.Expr) *ast.Ident {
	switch x := expr.(type) {
	case *ast.Ident:
		return x
	case *ast.StarExpr:
		return ident(x.X)
	case *ast.SelectorExpr:
		return ident(x.X)
	default:
		return nil
	}
}
