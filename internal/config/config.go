package config

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const _configFileName = "codemark.yaml"

// ReadIn will read the config file from `path` and return the `section` if any
// found. If no section is found then an empty map will be returned. If path is
// empty it will try to find a config file under $PWD. If no config file is
// found at all then
func ReadIn(path, section string) (map[string]any, error) {
	precedenceOrder := []string{path}
	defaultOrder, err := defaultPrecedenceOrder()
	if err != nil {
		return nil, err
	}
	precedenceOrder = append(precedenceOrder, defaultOrder...)
	var config map[string]any
	for _, path := range precedenceOrder {
		file, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			continue
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
	return config, nil
}

func defaultPrecedenceOrder() ([]string, error) {
	order := []string{}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	order = append(order, filepath.Join(wd, _configFileName))
	return order, nil
}
