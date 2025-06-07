package codemark

import (
	"go/ast"
)

func isEmbedded(field *ast.Field) bool {
	return len(field.Names) == 0
}

func isMethod(fn *ast.FuncDecl) bool {
	return fn.Recv != nil
}
