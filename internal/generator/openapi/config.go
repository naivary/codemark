package openapi

import (
	"github.com/goccy/go-yaml"

	"github.com/naivary/codemark/optionutil"
)

func newConfig(cfg map[string]any) (*config, error) {
	c := config{
		Schema: schemaConfig{
			Draft:     Draft("https://json-schema.org/draft/2020-12/schema"),
			IDBaseURL: "http://codemark.io/schemas",
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
	err = yaml.Unmarshal(data, &c)
	return &c, err
}

type config struct {
	Schema schemaConfig `yaml:"schema"`
}

type schemaConfig struct {
	Draft Draft `yaml:"draft"`
	// IDBaseURL is the base of the url controlled by you. It follows the following
	// url convention: http://<domain-controlled-by-you>/schemas
	IDBaseURL string `yaml:"idBaseURL"`

	Formats schemaFormats `yaml:"formats"`
}

type schemaFormats struct {
	Property Format `yaml:"property"`
	Filename Format `yaml:"filename"`
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
