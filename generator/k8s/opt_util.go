package k8s

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/naivary/codemark/api/doc"
	optionapi "github.com/naivary/codemark/api/option"
	"github.com/naivary/codemark/option"
	"github.com/naivary/codemark/registry"
)

const _typeName = ""

type Docer[T any] interface {
	Doc() T
}

func newRegistry() (registry.Registry, error) {
	reg := registry.InMemory()
	opts := slices.Concat(
		configMapOpts(),
		objectMetaOpts(),
		podOpts(),
		rbacOpts(),
		serviceAccountOpts(),
	)
	for _, opt := range opts {
		if err := reg.Define(opt); err != nil {
			return nil, err
		}
	}
	return reg, nil
}

func mustMakeOpt(name string, output any, isUnique bool, targets ...optionapi.Target) *optionapi.Option {
	rtype := reflect.TypeOf(output)
	doc := output.(Docer[doc.Option]).Doc()
	// undefined ident is needed to pass the validation of the name for the
	// MustMake function.
	if name == "" {
		name = strings.ToLower(rtype.Name())
	}
	ident := fmt.Sprintf("k8s:undefined:%s", name)
	opt := option.MustMake(ident, rtype, &doc, isUnique, targets...)
	return &opt
}

func makeOpts(resource string, opts ...*optionapi.Option) []*optionapi.Option {
	for _, opt := range opts {
		opt.Ident = fmt.Sprintf("k8s:%s:%s", resource, option.OptionOf(opt.Ident))
		if isResource(opt.Ident, "undefined") {
			opt.Ident = fmt.Sprintf("k8s:%s:%s", resource, opt.Output.Name())
		}
	}
	return opts
}

func isResource(ident, resource string) bool {
	return option.ResourceOf(ident) == resource
}
