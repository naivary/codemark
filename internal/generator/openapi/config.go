package openapi

import (
	"github.com/goccy/go-yaml"

	"github.com/naivary/codemark/optionutil"
)

func newConfig(cfg map[string]any) (*config, error) {
	c := config{
		Schema: schemaConfig{
			Draft:     "https://json-schema.org/draft/2020-12/schema",
			IDBaseURL: "",
			Formats: schemaFormats{
				Property: CamelCase,
				Filename: SnakeCase,
			},
		},
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	// err = validateConfig(data)
	// if err != nil {
	// 	return nil, err
	// }
	err = yaml.Unmarshal(data, &c)
	return &c, err
}

// func validateConfig(data []byte) error {
// 	const schemaFilePath = "./schemas/config.json"
// 	json, err := yaml.YAMLToJSON(data)
// 	if err != nil {
// 		return err
// 	}
// 	r := bytes.NewReader(json)
// 	c := jsonschema.NewCompiler()
// 	schm, err := c.Compile(schemaFilePath)
// 	if err != nil {
// 		return err
// 	}
// 	inst, err := jsonschema.UnmarshalJSON(r)
// 	if err != nil {
// 		return err
// 	}
// 	return schm.Validate(inst)
// }

// +openapi:schema:description="config options for the openapi generator"
type config struct {
	Schema schemaConfig `yaml:"schema"`
}

// +openapi:schema:description="config options for the schema model of openapi"
type schemaConfig struct {
	// +openapi:schema:enum=["https://json-schema.org/draft/2020-12/schema"]
	Draft string `yaml:"draft"`

	IDBaseURL string `yaml:"idBaseURL"`

	Formats schemaFormats `yaml:"formats"`
}

// +openapi:schema:description="available formats for the property and filename"
type schemaFormats struct {
	// +openapi:schema:enum=["snake_case", "camelCase", "pascalCase", "kebab-case", "ENV"]
	Property NamingConvention `yaml:"property"`
	// +openapi:schema:enum=["snake_case"]
	Filename NamingConvention `yaml:"filename"`
}

// Get is returning the default value of the identifier. If no default value is
// available nil will be returned.
func (c *config) Get(ident string) any {
	resource := optionutil.ResourceOf(ident)
	option := optionutil.OptionOf(ident)
	switch resource {
	case _schemaResource:
		return c.getSchema(option)
	default:
		return nil
	}
}

func (c *config) getSchema(opt string) any {
	switch opt {
	default:
		return nil
	}
}
