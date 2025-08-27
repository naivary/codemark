package filesystem

import (
	genv1 "github.com/naivary/codemark/api/generator/v1"
	outv1 "github.com/naivary/codemark/api/outputer/v1"
)

var _ outv1.Outputer = (*outputer)(nil)

type outputer struct{}

func (o *outputer) Name() string {
	return "fs"
}

func (o *outputer) Output(artifacts []*genv1.Artifact, cfg map[string]any) error {
	return nil
}
