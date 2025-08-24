package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

func ReadIn(path, section string) (map[string]any, error) {
	var config map[string]any
	file, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]any{}, nil
	}
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	cfg, isMap := config[section].(map[string]any)
	if !isMap {
		return make(map[string]any), nil
	}
	return cfg, nil
}
