package cmd

import (
	"github.com/goccy/go-yaml"

	"github.com/naivary/codemark/internal/config"
)

type cliConfig struct {
	DefaultOutputer string `yaml:"defaultOutputer"`
}

func newConfig(cfgFile string) (*cliConfig, error) {
	const configSection = "cli"
	c := cliConfig{
		DefaultOutputer: "fs",
	}
	cfg, err := config.ReadIn(cfgFile, configSection)
	if err != nil {
		return nil, err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &c)
	return &c, err
}
