package k8s

import "github.com/pelletier/go-toml/v2"

type config struct {
	Namespace  string `toml:"namespace"`
	NameFormat Format `toml:"name_format"`

	ConfigMap configMapConfig `toml:"configmap"`
}

type configMapConfig struct {
	KeyFormat Format `toml:"format"`
}

func newConfig(cfg map[string]any) (*config, error) {
	c := config{
		Namespace:  "default",
		NameFormat: SnakeCase,
		ConfigMap: configMapConfig{
			KeyFormat: CamelCase,
		},
	}
	data, err := toml.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(data, &c)
	return &c, err
}
