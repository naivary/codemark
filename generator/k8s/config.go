package k8s

import (
	"github.com/pelletier/go-toml/v2"

	"github.com/naivary/codemark/optionutil"
)

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

type config struct {
	Namespace  string `toml:"namespace"`
	NameFormat Format `toml:"name_format"`

	ConfigMap configMapConfig `toml:"configmap"`
}

func (c *config) Get(ident string) any {
	resource := optionutil.ResourceOf(ident)
	option := optionutil.OptionOf(ident)
	switch resource {
	case _configMapResource:
		return c.getConfigMap(option)
	}
	return nil
}

func (c *config) getConfigMap(option string) any {
	switch option {
	case "default":
		return nil
	case "immutable":
		return c.ConfigMap.Immutable
	case "format.key":
		return c.ConfigMap.KeyFormat
	default:
		return nil
	}
}

type configMapConfig struct {
	Immutable Immutable `toml:"immutable"`
	KeyFormat Format    `toml:"format"`
}
