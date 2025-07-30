package k8s

import (
	"bytes"
	"io"

	"github.com/goccy/go-yaml"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
)

func shouldGenerateConfigMap(strc *infov1.StructInfo) bool {
	for _, field := range strc.Fields {
		for ident := range field.Opts {
			if ident == "k8s:configmap:default" {
				return true
			}
		}
	}
	return false
}

func isMainFunc(fn infov1.FuncInfo) bool {
	return fn.Decl.Name.Name == "main"
}

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
