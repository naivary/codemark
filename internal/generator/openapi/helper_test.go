package openapi

import (
	"encoding/json"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	configer "github.com/naivary/codemark/internal/config"
	"github.com/naivary/codemark/loader"
)

func gen(path, cfgFile string) ([]*genv1.Artifact, error) {
	gen, err := New()
	if err != nil {
		return nil, err
	}
	proj, err := loader.Load(gen.Registry(), nil, path)
	if err != nil {
		return nil, err
	}
	if cfgFile == "" {
		return gen.Generate(proj, nil)
	}
	gensCfg, err := configer.ReadIn(path, "gens")
	if err != nil {
		return nil, err
	}
	cfg := gensCfg["openapi"].(map[string]any)
	return gen.Generate(proj, cfg)
}

// func gen(path string) ([]*genv1.Artifact, error) {
// 	mngr, err := generator.NewManager("")
// 	if err != nil {
// 		panic(err)
// 	}
// 	gen, err := New()
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = mngr.Add(gen)
// 	if err != nil {
// 		panic(err)
// 	}
// 	artifacts, err := mngr.Generate(nil, path)
// 	return artifacts[gen.Domain()], err
// }

func mustMarshal(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
