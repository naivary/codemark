package k8s

import (
	"bytes"

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
	for _, manifest := range manifests {
		if err := yaml.NewEncoder(&data).Encode(&manifest); err != nil {
			return nil, err
		}
		if _, err := data.WriteString("---"); err != nil {
			return nil, err
		}
	}
	return &genv1.Artifact{
		Name: name,
		Data: &data,
	}, nil
}
