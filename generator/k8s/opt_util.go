package k8s

import (
	"fmt"
	"reflect"
	"strings"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
	"github.com/naivary/codemark/optionutil"
)

const _typeName = ""

const (
	_unique    = true
	_repetable = false
)

type Docer[T any] interface {
	Doc() T
}

func mustMakeOpt(name string, output any, isUnique bool, targets ...optionv1.Target) *optionv1.Option {
	rtype := reflect.TypeOf(output)
	doc := output.(Docer[docv1.Option]).Doc()
	// undefined ident is needed to pass the validation of the name for the
	// MustMake function.
	if name == "" {
		name = strings.ToLower(rtype.Name())
	}
	ident := fmt.Sprintf("k8s:undefined:%s", name)
	opt := optionutil.MustMake(ident, rtype, &doc, isUnique, targets...)
	return &opt
}

func makeOpts(resource string, opts ...*optionv1.Option) []*optionv1.Option {
	for _, opt := range opts {
		opt.Ident = fmt.Sprintf("k8s:%s:%s", resource, optionutil.OptionOf(opt.Ident))
		if isResource(opt.Ident, "undefined") {
			opt.Ident = fmt.Sprintf("k8s:%s:%s", resource, opt.Output.Name())
		}
	}
	return opts
}

func isResource(ident, resource string) bool {
	return optionutil.ResourceOf(ident) == resource
}
