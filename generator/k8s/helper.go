package k8s

import (
	"slices"

	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/registry"
)

func newRegistry() (registry.Registry, error) {
	reg := registry.InMemory()
	defs := slices.Concat(configMapDefs(), objectMetaDefs())
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
			return nil, err
		}
	}
	return reg, nil
}

func shouldGenerateConfigMap(strc *loaderapi.StructInfo) bool {
	for _, field := range strc.Fields {
		for ident := range field.Defs {
			if ident == "k8s:configmap:default" {
				return true
			}
		}
	}
	return false
}
