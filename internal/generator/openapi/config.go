package openapi

import (
	"fmt"
	"net/url"

	"github.com/goccy/go-yaml"

	"github.com/naivary/codemark/optionutil"
)

func newConfig(cfg map[string]any) (*config, error) {
	c := config{
		Schema: schemaConfig{
			Draft:     "https://json-schema.org/draft/2020-12/schema",
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
	if err != nil {
		return nil, err
	}
	return &c, c.validate()
}

type config struct {
	Schema schemaConfig `yaml:"schema"`
}

type schemaConfig struct {
	Draft string `yaml:"draft"`
	// IDBaseURL is the base of the url controlled by you. It follows the following
	// url convention: http://<domain-controlled-by-you>/schemas
	IDBaseURL string `yaml:"idBaseURL"`

	Formats schemaFormats `yaml:"formats"`
}

type schemaFormats struct {
	Property NamingConvention `yaml:"property"`
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

func (c *config) validate() error {
	if _, err := url.ParseRequestURI(c.Schema.Draft); err != nil {
		return fmt.Errorf("the draft you have choosen is not a proper URL: %s", c.Schema.Draft)
	}
	return nil
}
