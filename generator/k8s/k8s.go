package k8s

import (
	"slices"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	"github.com/naivary/codemark/registry"
)

func newRegistry() (registry.Registry, error) {
	reg := registry.InMemory()
	opts := slices.Concat(
		configMapOpts(),
		objectMetaOpts(),
		rbacOpts(),
		serviceAccountOpts(),
	)
	for _, opt := range opts {
		if err := reg.Define(opt); err != nil {
			return nil, err
		}
	}
	return reg, nil
}

var _ genv1.Generator = (*k8sGenerator)(nil)

type k8sGenerator struct {
	reg registry.Registry
}

func New() (genv1.Generator, error) {
	reg, err := newRegistry()
	if err != nil {
		return nil, err
	}
	gen := &k8sGenerator{
		reg: reg,
	}
	return gen, nil
}

func (g k8sGenerator) Domain() string {
	return "k8s"
}

func (g k8sGenerator) Explain(ident string) string {
	_, err := g.reg.Get(ident)
	if err != nil {
		return ""
	}
	return ""
}

func (g k8sGenerator) Registry() registry.Registry {
	return g.reg
}

func (g k8sGenerator) Generate(proj infov1.Project, config map[string]any) ([]*genv1.Artifact, error) {
	artifacts := make([]*genv1.Artifact, 0, len(proj))
	for _, info := range proj {
		for _, strc := range info.Structs {
			if shouldGenerateConfigMap(strc) {
				cm, err := createConfigMap(strc)
				if err != nil {
					return nil, err
				}
				artifacts = append(artifacts, cm)
			}
		}
		for _, fn := range info.Funcs {
			if isMainFunc(fn) {
				// createPod(fn)
				rbac, err := createRBAC(fn)
				if err != nil {
					return nil, err
				}
				artifacts = append(artifacts, rbac)
			}
		}
	}
	return artifacts, nil
}
