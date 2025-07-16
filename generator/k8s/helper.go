package k8s

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/naivary/codemark/api"
	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/maker"
	"github.com/naivary/codemark/registry"
)

type docer interface {
	Doc() api.OptionDoc
}

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

func makeDefs(resource string, opts ...any) []*api.Definition {
	defs := make([]*api.Definition, 0, len(opts))
	for _, opt := range opts {
		to := reflect.TypeOf(opt)
		name := strings.ToLower(to.Name())
		ident := fmt.Sprintf("k8s:%s:%s", resource, name)
		doc := opt.(docer).Doc()
		def := maker.MustMakeDefWithDoc(ident, to, doc, doc.Targets...)
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
