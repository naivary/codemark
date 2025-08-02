package k8s

import (
	"maps"
	"reflect"
	"slices"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	"github.com/naivary/codemark/registry"
)

func newRegistry(resources ...Resourcer) (registry.Registry, error) {
	reg := registry.InMemory()
	for _, resource := range resources {
		opts := resource.Options()
		for _, opt := range opts {
			if err := reg.Define(opt); err != nil {
				return nil, err
			}
		}

	}
	for _, opt := range objectMetaOpts() {
		if err := reg.Define(opt); err != nil {
			return nil, err
		}
	}
	return reg, nil
}

var _ genv1.Generator = (*k8sGenerator)(nil)

type k8sGenerator struct {
	domain string

	reg registry.Registry

	resources map[reflect.Type][]Resourcer
}

func New() (genv1.Generator, error) {
	gen := &k8sGenerator{
		domain: "k8s",
		resources: map[reflect.Type][]Resourcer{
			reflect.TypeFor[*infov1.StructInfo](): {NewConfigMapResourcer()},
			reflect.TypeFor[*infov1.FuncInfo]():   {NewRBACResourcer()},
		},
	}
	resources := flatten(slices.Collect(maps.Values(gen.resources)))
	reg, err := newRegistry(resources...)
	if err != nil {
		return nil, err
	}
	gen.reg = reg
	return gen, nil
}

func (g k8sGenerator) Domain() string {
	return g.domain
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
	cfg, err := newConfig(config)
	if err != nil {
		return nil, err
	}
	for pkg, pkgInfo := range proj {
		infos := collectInfos(pkgInfo)
		for obj, info := range infos {
			metadata, err := createObjectMeta(info, cfg)
			if err != nil {
				return nil, err
			}
			infoType := reflect.TypeOf(info)
			for _, resource := range g.resources[infoType] {
				if !resource.CanCreate(info) {
					continue
				}
				artifact, err := resource.Create(pkg, obj, info, metadata, cfg)
				if err != nil {
					return nil, err
				}
				artifacts = append(artifacts, artifact)
			}

		}
	}
	return artifacts, nil
}
