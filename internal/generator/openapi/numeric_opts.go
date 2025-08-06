package openapi

import (
	"errors"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

type MultipleOf int64

func (m MultipleOf) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Multiple of",
		Default: "",
	}
}

func (m MultipleOf) apply(schema *Schema) error {
	i := int64(m)
	if i == 0 {
		return errors.New("a multiple of 0 is not allowed")
	}
	if schema.Type != integerType && schema.Type != numberType {
		return errors.New("multipleOf marker can only be used for integer-like types")
	}
	schema.MultipleOf = i
	return nil
}

type Minimum int64

func (m Minimum) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Minimum allowable number (non-exclusive)",
		Default: "",
	}
}

func (m Minimum) apply(schema *Schema) error {
	if schema.Type != numberType && schema.Type != integerType {
		return errors.New("minimum marker can only be applied to int types")
	}
	schema.Minimum = int64(m)
	return nil
}

type Maximum int64

func (m Maximum) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Maximum allowable number (non-exclusive)",
		Default: "",
	}
}

func (m Maximum) apply(schema *Schema) error {
	if schema.Type != numberType && schema.Type != integerType {
		return errors.New("maximum marker can only be applied to int types")
	}
	schema.Maximum = int64(m)
	return nil
}

type ExclusiveMinimum int64

func (e ExclusiveMinimum) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Exclusive minimum",
		Default: "",
	}
}

func (e ExclusiveMaximum) apply(schema *Schema) error {
	i := int64(e)
	if schema.Type != integerType && schema.Type != numberType {
		return errors.New("exclusiveMaximum marker can only be used for integer-like types")
	}
	schema.ExclusiveMaximum = i
	return nil
}

type ExclusiveMaximum int

func (e ExclusiveMaximum) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Exclusive maximum",
		Default: "",
	}
}

func (e ExclusiveMinimum) apply(schema *Schema) error {
	i := int64(e)
	if schema.Type != integerType && schema.Type != numberType {
		return errors.New("exclusiveMinimum marker can only be used for integer-like types")
	}
	schema.ExclusiveMinimum = i
	return nil
}
