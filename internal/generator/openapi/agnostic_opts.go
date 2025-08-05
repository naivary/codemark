package openapi

import (
	"errors"
	"net/url"

	"github.com/iancoleman/strcase"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

type Description string

func (d Description) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Description",
		Default: "",
	}
}

func (d Description) apply(schema *Schema) error {
	desc := string(d)
	if len(desc) == 0 {
		return errors.New("description for schema cannot be empty")
	}
	schema.Desc = desc
	return nil
}

type ID string

func (i ID) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "ID",
		Default: "",
	}
}

func (i ID) apply(schema *Schema, baseURL string, format Format) error {
	id := format.Format(string(i)) + ".json"
	idURL, err := url.JoinPath(baseURL, id)
	if err != nil {
		return err
	}
	schema.ID = idURL
	return nil
}

type Draft string

func (d Draft) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Draft",
		Default: "",
	}
}

func (d Draft) apply(schema *Schema) error {
	draft := string(d)
	if len(draft) == 0 {
		return errors.New("draft for schema cannot be empty")
	}
	// is valid url?
	if _, err := url.Parse(draft); err != nil {
		return err
	}
	schema.Draft = draft
	return nil
}

type Format string

const (
	SnakeCase  Format = "snake_case"
	CamelCase  Format = "camelCase"
	PascalCase Format = "PascalCase"
	KebabCase  Format = "kebab-case"
	Env        Format = "ENV"
)

func (k Format) Format(key string) string {
	switch k {
	case SnakeCase:
		return strcase.ToSnake(key)
	case CamelCase:
		return strcase.ToLowerCamel(key)
	case PascalCase:
		return strcase.ToCamel(key)
	case Env:
		return strcase.ToScreamingDelimited(key, '_', "", true)
	case KebabCase:
		return strcase.ToKebab(key)
	default:
		return ""
	}
}

func (k Format) Doc() docv1.Option {
	return docv1.Option{
		Desc:    `Format of the key. It be used to manipulate based on the context of the configuration. For example if the configuration is ssettable via environment variable it is useful to choose the env formation.`,
		Default: "camelCase",
	}
}
