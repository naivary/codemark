package openapi

import (
	"errors"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
)

type Required bool

func (r Required) Doc() docv1.Option {
	return docv1.Option{}
}

func (r Required) apply(schema *Schema, finfo *infov1.FieldInfo, cfg *config) error {
	if schema.Type != objectType {
		return errors.New("required can only be applied to objects")
	}
	fieldName := cfg.Schema.Formats.Property.Format(finfo.Ident.Name)
	schema.Required = append(schema.Required, fieldName)
	return nil
}

type DependentRequired []string

func (dr DependentRequired) Doc() docv1.Option {
	return docv1.Option{}
}

func (dr DependentRequired) apply(schema *Schema, cfg *config, finfo *infov1.FieldInfo) error {
	if schema.Type != objectType {
		return errors.New("dependentRequired can only be applied to objects")
	}
	fieldName := cfg.Schema.Formats.Property.Format(finfo.Ident.Name)
	for _, required := range dr {
		field := cfg.Schema.Formats.Property.Format(required)
		schema.DependentRequired[fieldName] = append(schema.DependentRequired[fieldName], field)
	}
	return nil
}
