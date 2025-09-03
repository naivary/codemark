//go:generate codemark gen -o openapi:fs ./... -- --fs.path=schemas
package openapi

import (
	"maps"
	"reflect"
	"slices"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	regv1 "github.com/naivary/codemark/api/registry/v1"
	"github.com/naivary/codemark/registry"
)

func newRegistry(resources ...Resourcer) (regv1.Registry, error) {
	reg := registry.InMemory()
	for _, resource := range resources {
		opts := resource.Options()
		for _, opt := range opts {
			if err := reg.Define(opt); err != nil {
				return nil, err
			}
		}

	}
	return reg, nil
}

const _domain = "openapi"

var _ genv1.Generator = (*openAPIGenerator)(nil)

func New() (genv1.Generator, error) {
	gen := &openAPIGenerator{
		resources: map[reflect.Type][]Resourcer{
			reflect.TypeFor[*infov1.StructInfo](): {NewSchemaResourcer()},
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

type openAPIGenerator struct {
	reg regv1.Registry

	resources map[reflect.Type][]Resourcer
}

func (g *openAPIGenerator) Domain() docv1.Domain {
	return docv1.Domain{
		Name: _domain,
		Desc: "",
	}
}

func (g *openAPIGenerator) Explain(ident string) string {
	doc, err := g.reg.DocOf(ident)
	if err != nil {
		return "not found"
	}
	return doc.Desc
}

func (g *openAPIGenerator) Registry() regv1.Registry {
	return g.reg
}

func (g *openAPIGenerator) Generate(proj infov1.Project, config map[string]any) ([]*genv1.Artifact, error) {
	cfg, err := newConfig(config)
	if err != nil {
		return nil, err
	}
	artifacts := make([]*genv1.Artifact, 0, len(proj))
	for pkg, pkgInfo := range proj {
		infos := collectInfos(pkgInfo)
		for obj, info := range infos {
			infoType := reflect.TypeOf(info)
			for _, resource := range g.resources[infoType] {
				if !resource.CanCreate(info) {
					continue
				}
				artifact, err := resource.Create(pkg, obj, info, cfg)
				if err != nil {
					return nil, err
				}
				artifacts = append(artifacts, artifact)
			}

		}
	}
	return artifacts, nil
}
