package k8s

import (
	loaderv1 "github.com/naivary/codemark/api/loader/v1"
)

func shouldGenerateConfigMap(strc *loaderv1.StructInfo) bool {
	for _, field := range strc.Fields {
		for ident := range field.Opts {
			if ident == "k8s:configmap:default" {
				return true
			}
		}
	}
	return false
}

func isMainFunc(fn loaderv1.FuncInfo) bool {
	return fn.Decl.Name.Name == "main"
}
