package openapi

import (
	"encoding/json"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	"github.com/naivary/codemark/generator"
)

func gen(path string) ([]*genv1.Artifact, error) {
	mngr, err := generator.NewManager("codemark.yaml")
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
	artifacts, err := mngr.Generate(nil, path)
	return artifacts[gen.Domain()], err
}

func mustMarshal(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
