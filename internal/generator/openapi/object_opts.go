package openapi

import (
	"errors"
	"fmt"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
)

// TODO: Pattern Property should be implemented

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

type MutuallyExclusive []string

func (me MutuallyExclusive) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "MutuallyExclusive is allowing you to define which field are mutually exclusive to each other. Be careful because the generator cannot verify that the logic will make sense in a high degree. It will only validate simple logic contradictions",
		Default: "",
	}
}

func (me MutuallyExclusive) apply(root *Schema, finfo *infov1.FieldInfo, sinfo *infov1.StructInfo, cfg *config) error {
	if hasOpt(finfo, "openapi:schema:required") {
		return fmt.Errorf(
			`contradicting logic: mutuallyExclusive and required can never be 
			fullfilled by a schema. Remove the required marker: %s`,
			finfo.Ident.Name,
		)
	}
	for _, fieldName := range me {
		if !sinfo.HasField(fieldName) {
			return fmt.Errorf(
				"the field you are trying to reference in mutuallyExclusive does not exist in your struct: field `%s` in struct `%s`",
				fieldName,
				sinfo.Spec.Name.Name,
			)
		}
		refFieldInfo := sinfo.GetFieldByIdent(fieldName)
		if hasOpt(refFieldInfo, "openapi:schema:required") {
			return fmt.Errorf("the reference field in mutuallyExclusive cannot be marked required: %s", fieldName)
		}
	}
	root.OneOf = append(root.OneOf, me.newOneOfSchema(finfo.Ident.Name, cfg)...)
	return nil
}

func (me MutuallyExclusive) newOneOfSchema(field string, cfg *config) []*Schema {
	oneOf := []*Schema{}
	exclusiveFields := []string{cfg.Schema.Formats.Property.Format(field)}
	for _, e := range me {
		exclusiveFields = append(exclusiveFields, cfg.Schema.Formats.Property.Format(e))
	}
	for _, field := range exclusiveFields {
		s := &Schema{
			Required: []string{field},
			Not: &Schema{
				AnyOf: make([]*Schema, 0, len(exclusiveFields)-1),
			},
		}
		for _, exclusiveField := range exclusiveFields {
			if field == exclusiveField {
				continue
			}
			anyOf := &Schema{
				Required: []string{exclusiveField},
			}
			s.Not.AnyOf = append(s.Not.AnyOf, anyOf)
		}
		oneOf = append(oneOf, s)
	}
	return oneOf
}
