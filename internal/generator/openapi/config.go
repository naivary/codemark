package openapi

import (
	"bytes"
	"embed"
	"fmt"
	"net/url"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schemas/*
var schemaFs embed.FS

type embedLoader struct{}

func (e *embedLoader) Load(urll string) (any, error) {
	u, err := url.Parse(urll)
	if err != nil {
		return "", err
	}
	// Find where "schemas/" starts
	idx := strings.Index(u.Path, "schemas/")
	if idx == -1 {
		return "", fmt.Errorf(`"schemas/" not found in path`)
	}
	// Return everything from schemas/ onward
	urll = u.Path[idx:]
	data, err := schemaFs.ReadFile(urll)
	if err != nil {
		return nil, err
	}
	return jsonschema.UnmarshalJSON(bytes.NewReader(data))
}

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
	err = validateConfig(data)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &c)
	return &c, err
}

func validateConfig(data []byte) error {
	json, err := yaml.YAMLToJSON(data)
	if err != nil {
		return err
	}
	r := bytes.NewReader(json)
	c := jsonschema.NewCompiler()
	c.UseLoader(jsonschema.SchemeURLLoader{
		"file": &embedLoader{},
	})
	const schemaFilePath = "schemas/config.json"
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
