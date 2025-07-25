package k8s

import (
	loaderapi "github.com/naivary/codemark/api/loader"
)

func shouldGenerateConfigMap(strc *loaderapi.StructInfo) bool {
	for _, field := range strc.Fields {
		for ident := range field.Opts {
			if ident == "k8s:configmap:default" {
				return true
			}
		}
	}
	return false
}

func isMainFunc(fn loaderapi.FuncInfo) bool {
	return fn.Decl.Name.Name == "main"
}
