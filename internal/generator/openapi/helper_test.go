package openapi

import (
	genv1 "github.com/naivary/codemark/api/generator/v1"
	"github.com/naivary/codemark/generator"
)

func gen(path string) ([]*genv1.Artifact, error) {
	mngr, err := generator.NewManager()
	if err != nil {
		panic(err)
	}
	gen, err := New()
	if err != nil {
		panic(err)
	}
	err = mngr.Add(gen)
	if err != nil {
		panic(err)
	}
	artifacts, err := mngr.Generate(path, gen.Domain())
	return artifacts[gen.Domain()], err
}
