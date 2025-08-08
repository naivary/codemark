package openapi

import (
	"errors"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

func isNumberOrInteger(typ jsonType) bool {
	return typ == numberType || typ == integerType
}

type MultipleOf int64

func (m MultipleOf) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Multiple of",
		Default: "",
	}
}

func (m MultipleOf) apply(schema *Schema) error {
	multipleOf := int64(m)
	if multipleOf == 0 {
		return errors.New("a multiple of 0 is not allowed")
	}
	if isNumberOrInteger(schema.Type) {
		schema.MultipleOf = multipleOf
		return nil
	}
	if isNumberOrInteger(schema.Items.Type) {
		schema.Items.MultipleOf = multipleOf
		return nil
	}
	if isNumberOrInteger(schema.AdditionalProperties.Type) {
		schema.AdditionalProperties.MultipleOf = multipleOf
		return nil
	}
	return errors.New(
		"multipleOf is only appliable to int[8,16,32,64]/float[32,64], objects with value of type int/float and arrays/slices of type int/float",
	)
}

type Minimum int64

func (m Minimum) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Minimum allowable number (non-exclusive)",
		Default: "",
	}
}

func (m Minimum) apply(schema *Schema) error {
	minimum := int64(m)
	if isNumberOrInteger(schema.Type) {
		schema.Minimum = minimum
		return nil
	}
	if isNumberOrInteger(schema.Items.Type) {
		schema.Items.Minimum = minimum
		return nil
	}
	if isNumberOrInteger(schema.AdditionalProperties.Type) {
		schema.AdditionalProperties.Minimum = minimum
		return nil
	}
	return errors.New(
		"minimum is only appliable to int[8,16,32,64]/float[32,64], objects with value of type int/float and arrays/slices of type int/float",
	)
}

type Maximum int64

func (m Maximum) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Maximum allowable number (non-exclusive)",
		Default: "",
	}
}

func (m Maximum) apply(schema *Schema) error {
	maximum := int64(m)
	if isNumberOrInteger(schema.Type) {
		schema.Maximum = maximum
		return nil
	}
	if isNumberOrInteger(schema.Items.Type) {
		schema.Items.Maximum = maximum
		return nil
	}
	if isNumberOrInteger(schema.AdditionalProperties.Type) {
		schema.AdditionalProperties.Maximum = maximum
		return nil
	}
	return errors.New(
		"maximum is only appliable to int[8,16,32,64]/float[32,64], objects with value of type int/float and arrays/slices of type int/float",
	)
}

type ExclusiveMinimum int64

func (e ExclusiveMinimum) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Exclusive minimum",
		Default: "",
	}
}

func (e ExclusiveMaximum) apply(schema *Schema) error {
	exclusiveMaximum := int64(e)
	if isNumberOrInteger(schema.Type) {
		schema.ExclusiveMaximum = exclusiveMaximum
		return nil
	}
	if isNumberOrInteger(schema.Items.Type) {
		schema.Items.ExclusiveMaximum = exclusiveMaximum
		return nil
	}
	if isNumberOrInteger(schema.AdditionalProperties.Type) {
		schema.AdditionalProperties.ExclusiveMaximum = exclusiveMaximum
		return nil
	}
	return errors.New(
		"exclusiveMaximum is only appliable to int[8,16,32,64]/float[32,64], objects with value of type int/float and arrays/slices of type int/float",
	)
}

type ExclusiveMaximum int

func (e ExclusiveMaximum) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Exclusive maximum",
		Default: "",
	}
}

func (e ExclusiveMinimum) apply(schema *Schema) error {
	exclusiveMinimum := int64(e)
	if isNumberOrInteger(schema.Type) {
		schema.ExclusiveMinimum = exclusiveMinimum
		return nil
	}
	if isNumberOrInteger(schema.Items.Type) {
		schema.Items.ExclusiveMinimum = exclusiveMinimum
		return nil
	}
	if isNumberOrInteger(schema.AdditionalProperties.Type) {
		schema.AdditionalProperties.ExclusiveMinimum = exclusiveMinimum
		return nil
	}
	return errors.New(
		"exclusiveMinimum is only appliable to int[8,16,32,64]/float[32,64], objects with value of type int/float and arrays/slices of type int/float",
	)
}
