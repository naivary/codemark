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
		podOpts(),
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

func (g generator) Generate(proj infov1.Project) ([]*genv1.Artifact, error) {
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
