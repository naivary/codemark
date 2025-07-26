package k8s

import (
	"bytes"
	"encoding/json"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	loaderv1 "github.com/naivary/codemark/api/loader/v1"
	"github.com/naivary/codemark/registry"
)

var _ genv1.Generator = (*generator)(nil)

type generator struct {
	reg registry.Registry
}

func New() (genv1.Generator, error) {
	reg, err := newRegistry()
	if err != nil {
		return nil, err
	}
	gen := &generator{
		reg: reg,
	}
	return gen, nil
}

func (g generator) Domain() string {
	return "k8s"
}

func (g generator) Explain(ident string) string {
	option, err := g.reg.Get(ident)
	if err != nil {
		return ""
	}
	_ = option
	return ""
}

func (g generator) Registry() registry.Registry {
	return g.reg
}

func (g generator) Generate(proj loaderv1.Project) ([]genv1.Artifact, error) {
	artifacts := make([]genv1.Artifact, 0, len(proj))
	for _, info := range proj {
		for _, strc := range info.Structs {
			if shouldGenerateConfigMap(strc) {
				cm, err := createConfigMap(strc)
				if err != nil {
					return nil, err
				}
				artifacts = append(artifacts, newArtifact(cm.Name, cm))
			}
		}
		for _, fn := range info.Funcs {
			if isMainFunc(fn) {
				pod, err := createPod(fn)
				if err != nil {
					return nil, err
				}
				role, err := createRBACRole(fn)
				if err != nil {
					return nil, err
				}
				artifacts = append(artifacts, newArtifact(pod.Name, pod))
				artifacts = append(artifacts, newArtifact(role.Name, role))
			}
		}
	}
	return artifacts, nil
}

func newArtifact(name string, data any) genv1.Artifact {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(data)
	if err != nil {
		panic(err)
	}
	return genv1.Artifact{
		Name:        name,
		ContentType: "application/json",
		Data:        &b,
	}
}
