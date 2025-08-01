package k8s

import (
	"bytes"

	"github.com/goccy/go-yaml"
	"github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/naivary/codemark/optionutil"
)

func validateConfig(data []byte) error {
	const schemaFilePath = "./schema/k8s.json"
	json, err := yaml.YAMLToJSON(data)
	if err != nil {
		return err
	}
	r := bytes.NewReader(json)
	c := jsonschema.NewCompiler()
	schm, err := c.Compile(schemaFilePath)
	if err != nil {
		return err
	}
	inst, err := jsonschema.UnmarshalJSON(r)
	if err != nil {
		return err
	}
	return schm.Validate(inst)
}

func newConfig(cfg map[string]any) (*config, error) {
	c := config{
		Namespace:  Namespace("default"),
		NameFormat: KebabCase,
		ConfigMap: configMapConfig{
			KeyFormat: CamelCase,
		},
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	err = validateConfig(data)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &c)
	return &c, err
}

type configMapConfig struct {
	KeyFormat Format    `yaml:"key_format"`
	immutable Immutable `yaml:"immutable"`
}

type config struct {
	Namespace  Namespace `yaml:"namespace"`
	NameFormat Format    `yaml:"name_format"`

	ConfigMap configMapConfig `yaml:"configMap"`
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
