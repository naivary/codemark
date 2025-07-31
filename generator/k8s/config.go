package k8s

import (
	"github.com/pelletier/go-toml/v2"

	"github.com/naivary/codemark/optionutil"
)

func newConfig(cfg map[string]any) (*config, error) {
	c := config{
		Namespace:  Namespace("default"),
		NameFormat: SnakeCase,
		ConfigMap: configMapConfig{
			KeyFormat: CamelCase,
		},
	}
	data, err := toml.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	return &c, toml.Unmarshal(data, &c)
}

type configMapConfig struct {
	Immutable Immutable `toml:"immutable"`
	KeyFormat Format    `toml:"format"`
}

type config struct {
	Namespace  Namespace `toml:"namespace"`
	NameFormat Format    `toml:"name_format"`

	ConfigMap configMapConfig `toml:"configmap"`
}

func (c *config) Get(ident string) any {
	resource := optionutil.ResourceOf(ident)
	option := optionutil.OptionOf(ident)
	switch resource {
	case _configMapResource:
		return c.getConfigMap(option)
	case _objectMetaResource:
		return c.getObjectMeta(option)
	default:
		return nil
	}
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

func (c *config) getObjectMeta(option string) any {
	switch option {
	case "namespace":
		return c.Namespace
	case "format.name":
		return c.NameFormat
	default:
		return nil
	}
}
