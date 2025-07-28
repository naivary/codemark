package k8s

import (
	infov1 "github.com/naivary/codemark/api/info/v1"
)

func shouldGenerateConfigMap(strc *infov1.StructInfo) bool {
	for _, field := range strc.Fields {
		for ident := range field.Opts {
			if ident == "k8s:configmap:default" {
				return true
			}
		}
	}
	return false
}

func isMainFunc(fn infov1.FuncInfo) bool {
	return fn.Decl.Name.Name == "main"
}
