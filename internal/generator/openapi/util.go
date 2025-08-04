package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
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
			opt.Ident = fmt.Sprintf("k8s:%s:%s", resource, opt.Type.Name())
		}
	}
	return opts
}

func isResource(ident, resource string) bool {
	return optionutil.ResourceOf(ident) == resource
}

func newArtifact(name string, manifests ...any) (*genv1.Artifact, error) {
	var data bytes.Buffer
	enc := json.NewEncoder(&data)
	for _, manifest := range manifests {
		if err := enc.Encode(&manifest); err != nil {
			return nil, err
		}
	}
	return &genv1.Artifact{
		Name: name,
		Data: &data,
	}, nil
}

func setDefaults(opts []*optionv1.Option, info infov1.Info, cfg *config, target optionv1.Target) {
	for _, opt := range opts {
		v := info.Options()[opt.Ident]
		if len(v) != 0 || !slices.Contains(opt.Targets, target) {
			continue
		}
		deflt := cfg.Get(opt.Ident)
		if deflt != nil {
			info.Options()[opt.Ident] = []any{deflt}
		}
	}
}
