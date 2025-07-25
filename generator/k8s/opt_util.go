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

type Docer[T any] interface {
	Doc() T
}

type optionMakeParams struct {
	typ      any
	targets  []optionapi.Target
	isUnique bool
	name     string
}

func newOption(typ any, isUnique bool, targets ...optionapi.Target) optionMakeParams {
	return optionMakeParams{
		typ:      typ,
		isUnique: isUnique,
		targets:  targets,
	}
}

func newRegistry() (registry.Registry, error) {
	reg := registry.InMemory()
	opts := slices.Concat(
		configMapOpts(),
		objectMetaOpts(),
		podOpts(),
		rbacOpts(),
	)
	for _, opt := range opts {
		if err := reg.Define(opt); err != nil {
			return nil, err
		}
	}
	return reg, nil
}

func makeOpts(resource string, optionTypes ...optionMakeParams) []*optionapi.Option {
	opts := make([]*optionapi.Option, 0, len(optionTypes))
	for _, p := range optionTypes {
		to := reflect.TypeOf(p.typ)
		name := strings.ToLower(to.Name())
		ident := fmt.Sprintf("k8s:%s:%s", resource, name)
		doc := p.typ.(Docer[doc.Option]).Doc()
		opt := option.MustMakeWithDoc(ident, to, doc, p.targets...)
		opt.IsUnique = p.isUnique
		opts = append(opts, opt)
	}
	return opts
}

func isResource(ident, resource string) bool {
	s := strings.Split(ident, ":")
	if len(s) != 3 {
		return false
	}
	return s[1] == resource
}
