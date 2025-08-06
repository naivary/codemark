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
