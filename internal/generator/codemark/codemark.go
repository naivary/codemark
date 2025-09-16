package codemark

import (
	"slices"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	regv1 "github.com/naivary/codemark/api/registry/v1"
	"github.com/naivary/codemark/registry"
)

const _domain = "codemark"

var _ genv1.Generator = (*codemarkGenerator)(nil)

type codemarkGenerator struct {
	reg regv1.Registry

	optDocRes *optDocResourcer

	configDocRes *configDocResourcer
}

func New() (genv1.Generator, error) {
	gen := &codemarkGenerator{
		optDocRes:    NewOptDocResourcer(),
		configDocRes: NewConfigDocResourcer(),
	}
	reg, err := gen.newRegistry()
	if err != nil {
		return nil, err
	}
	gen.reg = reg
	return gen, nil
}

func (c codemarkGenerator) Domain() docv1.Domain {
	return docv1.Domain{
		Name: _domain,
		Desc: `Generator to make your life easier to develop new generators and kinda prove the point of this whole project :)`,
	}
}

func (c *codemarkGenerator) Resources() map[string]*docv1.Resource {
	return nil
}

func (c *codemarkGenerator) Registry() regv1.Registry {
	return c.reg
}

func (c *codemarkGenerator) ConfigDoc() map[string]docv1.Config {
	return nil
}

func (c *codemarkGenerator) Generate(proj infov1.Project, config map[string]any) ([]*genv1.Artifact, error) {
	artifacts := make([]*genv1.Artifact, 0)
	optDoc := NewOptDocResourcer()
	for pkg, info := range proj {
		for _, s := range info.Structs {
			c.configDocRes.Create(pkg, s, proj)
		}
	}
	for pkg, info := range proj {
		artifact, err := optDoc.Create(pkg, info.Named)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts, nil
}

func (c codemarkGenerator) newRegistry() (regv1.Registry, error) {
	reg := registry.InMemory()
	opts := slices.Concat(
		c.optDocRes.Options(),
		c.configDocRes.Options(),
	)
	for _, opt := range opts {
		err := reg.Define(opt)
		if err != nil {
			return nil, err
		}
	}
	return reg, nil
}
