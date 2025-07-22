package k8s

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/naivary/codemark/api/core"
	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/option"
	"github.com/naivary/codemark/registry"
)

func newRegistry() (registry.Registry, error) {
	reg := registry.InMemory()
	opts := slices.Concat(configMapOpts(), objectMetaOpts(), podOpts())
	for _, opt := range opts {
		if err := reg.Define(opt); err != nil {
			return nil, err
		}
	}
	return reg, nil
}

func makeDefs(resource string, optionTypes map[any][]core.Target) []*core.Option {
	opts := make([]*core.Option, 0, len(optionTypes))
	for opt, targets := range optionTypes {
		to := reflect.TypeOf(opt)
		name := strings.ToLower(to.Name())
		ident := fmt.Sprintf("k8s:%s:%s", resource, name)
		doc := opt.(core.Docer[core.OptionDoc]).Doc()
		opt := option.MustMakeWithDoc(ident, to, doc, targets...)
		opts = append(opts, opt)
	}
	return opts
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
		for ident := range field.Opts {
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
