package openapi

import (
	"errors"
	"fmt"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
)

// TODO: Pattern Property should be implemented und AdditionalProperties

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

func (dr DependentRequired) apply(root *Schema, cfg *config, finfo *infov1.FieldInfo, structInfo *infov1.StructInfo) error {
	if root.Type != objectType {
		return errors.New("dependentRequired can only be applied to objects")
	}
	fieldName := cfg.Schema.Formats.Property.Format(finfo.Ident.Name)
	for _, required := range dr {
		if !structInfo.HasField(required) {
			return fmt.Errorf(
				"the field you are trying to reference in dependentRequired does not exist in your struct: field `%s` in struct `%s`",
				required,
				structInfo.Spec.Name.Name,
			)
		}
		field := cfg.Schema.Formats.Property.Format(required)
		root.DependentRequired[fieldName] = append(root.DependentRequired[fieldName], field)
	}
	return nil
}
