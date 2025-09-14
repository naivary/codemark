package codemark

import (
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
	cfg := gensCfg["codemark"].(map[string]any)
	return gen.Generate(proj, cfg)
}
