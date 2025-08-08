package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/types"
	"reflect"
	"slices"

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
	format := CamelCase
	rtype := reflect.TypeOf(output)
	doc := output.(Docer[docv1.Option]).Doc()
	if name == "" {
		name = format.Format(rtype.Name())
	}
	// undefined ident is needed to pass the validation of the name for the
	// MustMake function.
	ident := fmt.Sprintf("%s:undefined:%s", _domain, name)
	opt := optionutil.MustMake(ident, rtype, &doc, isUnique, targets...)
	return &opt
}

func makeOpts(resource string, opts ...*optionv1.Option) []*optionv1.Option {
	for _, opt := range opts {
		opt.Ident = fmt.Sprintf("%s:%s:%s", _domain, resource, optionutil.OptionOf(opt.Ident))
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
	artifact := &genv1.Artifact{
		Name: name,
		Data: bytes.NewBuffer(nil),
	}
	for _, manifest := range manifests {
		err := json.NewEncoder(artifact.Data).Encode(&manifest)
		if err != nil {
			return nil, err
		}
	}
	return artifact, nil
}

func setDefaults(opts []*optionv1.Option, info infov1.Info, cfg *config, target optionv1.Target, defaults map[string]any) {
	for _, opt := range opts {
		v := info.Options()[opt.Ident]
		if len(v) != 0 || !slices.Contains(opt.Targets, target) {
			continue
		}
		// TODO: multi value defaults are not supported rn
		deflt := cfg.Get(opt.Ident)
		if deflt == nil && defaults != nil {
			deflt = defaults[opt.Ident]
		}
		if deflt != nil {
			info.Options()[opt.Ident] = []any{deflt}
		}
	}
}

func flatten[T any](lists [][]T) []T {
	var res []T
	for _, list := range lists {
		res = append(res, list...)
	}
	return res
}

func collectInfos(info *infov1.Information) map[types.Object]infov1.Info {
	all := make(map[types.Object]infov1.Info, infoCap(info))
	for obj, str := range info.Structs {
		all[obj] = str
	}
	for obj, alias := range info.Aliases {
		all[obj] = alias
	}
	for obj, v := range info.Vars {
		all[obj] = v
	}
	for obj, c := range info.Consts {
		all[obj] = c
	}
	for obj, fn := range info.Funcs {
		all[obj] = fn
	}
	for obj, iface := range info.Ifaces {
		all[obj] = iface
	}
	for obj, imp := range info.Imports {
		all[obj] = imp
	}
	for obj, named := range info.Named {
		all[obj] = named
	}
	// TODO: Need types.Object for file or a fake one atleast
	// for filename, file := range info.Files {
	// 	all[nil] = file
	// }
	return all
}

func infoCap(info *infov1.Information) int {
	return len(
		info.Structs,
	) + len(
		info.Aliases,
	) + len(
		info.Vars,
	) + len(
		info.Consts,
	) + len(
		info.Funcs,
	) + len(
		info.Ifaces,
	) + len(
		info.Imports,
	) + len(
		info.Named,
	) + len(
		info.Files,
	)
}

func hasOpt(info infov1.Info, ident string) bool {
	v := info.Options()[ident]
	return len(v) != 0
}
