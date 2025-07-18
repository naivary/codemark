package k8s

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/naivary/codemark/api/core"
	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/maker"
	"github.com/naivary/codemark/registry"
)

func newRegistry() (registry.Registry, error) {
	reg := registry.InMemory()
	defs := slices.Concat(configMapOpts(), objectMetaOpts(), podOpts())
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
			return nil, err
		}
	}
	return reg, nil
}

func makeDefs(resource string, opts map[any][]core.Target) []*core.Option {
	defs := make([]*core.Option, 0, len(opts))
	for opt, targets := range opts {
		to := reflect.TypeOf(opt)
		name := strings.ToLower(to.Name())
		ident := fmt.Sprintf("k8s:%s:%s", resource, name)
		doc := opt.(core.Docer[core.OptionDoc]).Doc()
		def := maker.MustMakeOptWithDoc(ident, to, doc, targets...)
		defs = append(defs, def)
	}
	return defs
}

func keys[K comparable, V any](m map[K]V) []K {
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
