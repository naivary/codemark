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
		Desc: "Generate OpenAPI Specifications",
	}
}

func (g *openAPIGenerator) Resources() map[string]*docv1.Resource {
	return map[string]*docv1.Resource{
		_schemaResource: {Desc: "Generate an OpenAPI compatible JSON Schema"},
	}
}

func (g *openAPIGenerator) Registry() regv1.Registry {
	return g.reg
}

func (g *openAPIGenerator) ConfigDoc() map[string]docv1.Config {
	return map[string]docv1.Config{
		"schema": {
			Description: `Defines the full set of configuration options that control the generation of JSON Schemas compliant with the OpenAPI Specification. This includes structural metadata, type definitions, validation constraints, and any OpenAPI-specific extensions required for interoperability.`,
			Options: map[string]docv1.Config{
				"draft": {
					Default:     "https://json-schema.org/draft/2020-12/schema",
					Description: `Specifies the JSON Schema draft version to use for validation. Currently, this setting is not configurable — the default value is the only supported option. It exists for future compatibility and may allow other drafts in later releases.`,
				},
				"idBaseURL": {
					Default:     "",
					Description: `Sets the base URL for the $id field in generated JSON Schemas. The base URL is prepended to schema identifiers so they can be resolved consistently. If left empty (default), schemas will not have a network-resolvable $id and will only be referenceable from the local filesystem.`,
				},
				"formats": {
					Description: "Controls the output style of generated JSON Schemas. Use this option to align schema formatting (e.g., indentation, line wrapping, property ordering) with your organization’s conventions.",
					Options: map[string]docv1.Config{
						"property": {
							Default:     "camelCase",
							Description: "Defines how object properties are represented in the generated JSON Schema. This controls the formatting style (e.g. naming conventions) to ensure consistency with your organization’s schema style.",
						},
						"filename": {
							Default:     "snake_case",
							Description: "Specifies the naming convention for generated files. Use this option to control how schema filenames are produced (e.g., snake_case, kebab-case, PascalCase), ensuring they align with your project or organizational standards.",
						},
					},
				},
			},
		},
	}
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
