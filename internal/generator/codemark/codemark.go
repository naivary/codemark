package codemark

import (
	docv1 "github.com/naivary/codemark/api/doc/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	regv1 "github.com/naivary/codemark/api/registry/v1"
)

const _domain = "codemark"

var _ genv1.Generator = (*codemarkGenerator)(nil)

type codemarkGenerator struct {
	reg regv1.Registry
}

func New() genv1.Generator {
	return &codemarkGenerator{}
}

func (c codemarkGenerator) Domain() docv1.Domain {
	return docv1.Domain{
		Name: _domain,
		Desc: `Generator to make your life easier to develop new generators and kinda prove the point of this whole project :)`,
	}
}

func (c codemarkGenerator) Resources() map[string]*docv1.Resource {
	return nil
}

func (c codemarkGenerator) Registry() regv1.Registry {
	return nil
}

func (c codemarkGenerator) ConfigDoc() map[string]docv1.Config {
	return nil
}

func (c codemarkGenerator) Generate(proj infov1.Project, config map[string]any) ([]*genv1.Artifact, error) {
	artifacts := make([]*genv1.Artifact, 0)
	optDoc := NewOptDocResourcer()
	for pkg, info := range proj {
		artifact, err := optDoc.Create(pkg, info.Named)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts, nil
}
