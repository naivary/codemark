package k8s

import (
	"slices"

	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/registry"
)

func newRegistry() (registry.Registry, error) {
	reg := registry.InMemory()
	defs := slices.Concat(configMapDefs(), objectMetaDefs(), podDefs())
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
			return nil, err
		}
	}
	return reg, nil
}

func keysOfMap[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
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

func shouldGeneratePod(fn loaderapi.FuncInfo) bool {
	return fn.Decl.Name.Name == "main"
}
