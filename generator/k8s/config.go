package k8s

import (
	"errors"
	"fmt"

	"github.com/pelletier/go-toml/v2"

	"github.com/naivary/codemark/optionutil"
)

// TODO: validateConfig
func validateConfig(cfg *config) error {
	const formatTestString = "testString"
	if cfg.Namespace == "" {
		return errors.New("default namespace in codemark.toml cannot be empty")
	}
	if s := cfg.NameFormat.Format(formatTestString); s == "" {
		return fmt.Errorf("format not supported: %s", cfg.NameFormat)
	}
	if s := cfg.ConfigMap.KeyFormat.Format(formatTestString); s == "" {
		return fmt.Errorf("format not supported: %s", cfg.ConfigMap.KeyFormat)
	}
	return nil
}

func newConfig(cfg map[string]any) (*config, error) {
	c := config{
		Namespace:  Namespace("default"),
		NameFormat: KebabCase,
		ConfigMap: configMapConfig{
			KeyFormat: CamelCase,
		},
	}
	data, err := toml.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return &c, validateConfig(&c)
}

type configMapConfig struct {
	KeyFormat Format    `toml:"key_format"`
	immutable Immutable `toml:"immutable"`
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
		return c.ConfigMap.immutable
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
