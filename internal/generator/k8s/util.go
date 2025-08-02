package k8s

import (
	"bytes"
	"go/types"
	"io"
	"slices"

	"github.com/goccy/go-yaml"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

func newArtifact(name string, manifests ...any) (*genv1.Artifact, error) {
	var data bytes.Buffer
	enc := yaml.NewEncoder(&data)
	defer enc.Close()
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

func mergeArtifacts(artifacts ...*genv1.Artifact) (*genv1.Artifact, error) {
	if len(artifacts) == 0 {
		return nil, nil
	}
	base := artifacts[0]
	for _, artifact := range artifacts[1:] {
		_, err := io.WriteString(base.Data, "---\n")
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(base.Data, artifact.Data)
		if err != nil {
			return nil, err
		}
	}
	return base, nil
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
