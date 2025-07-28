package k8s

import (
	"fmt"
	"reflect"
	"slices"
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

var (
	_required any = nil
	// optional is something non nil signaling the underlying function used to
	// use the default of the type
	_optional any = new(any)
)

type Docer[T any] interface {
	Doc() T
}

// TODO: write test for this function
func setOptsDefaults(opts []*optionv1.Option, infoOpts map[string][]any, targets ...optionv1.Target) {
	for _, opt := range opts {
		if opt.IsRequired() {
			// requried options will be skipped
			continue
		}
		if !slices.Equal(opt.Targets, targets) && !slices.Contains(targets, optionv1.TargetAny) {
			continue
		}
		v := infoOpts[opt.Ident]
		if len(v) > 0 {
			continue
		}
		infoOpts[opt.Ident] = append(v, opt.Default)
	}
}

func mustMakeOpt(name string, output, defult any, isUnique bool, targets ...optionv1.Target) *optionv1.Option {
	rtype := reflect.TypeOf(output)
	doc := output.(Docer[docv1.Option]).Doc()
	// undefined ident is needed to pass the validation of the name for the
	// MustMake function.
	if name == "" {
		name = strings.ToLower(rtype.Name())
	}
	if defult != nil {
		defult = output
	}
	ident := fmt.Sprintf("k8s:undefined:%s", name)
	opt := optionutil.MustMake(ident, rtype, defult, &doc, isUnique, targets...)
	return &opt
}

func makeOpts(resource string, opts ...*optionv1.Option) []*optionv1.Option {
	for _, opt := range opts {
		opt.Ident = fmt.Sprintf("k8s:%s:%s", resource, optionutil.OptionOf(opt.Ident))
		if isResource(opt.Ident, "undefined") {
			opt.Ident = fmt.Sprintf("k8s:%s:%s", resource, opt.Type.Name())
		}
	}
	return opts
}

func isResource(ident, resource string) bool {
	return optionutil.ResourceOf(ident) == resource
}
