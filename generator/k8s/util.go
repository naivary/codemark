package k8s

import (
	"bytes"
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
		defaultt := cfg.Get(opt.Ident)
		if defaultt != nil {
			info.Options()[opt.Ident] = []any{defaultt}
		}
	}
}
