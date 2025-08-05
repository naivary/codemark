package openapi

import (
	"errors"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

type MultipleOf int

func (m MultipleOf) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Multiple of",
		Default: "",
	}
}

type Minimum int

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
	i := int64(m)
	schema.Minimum = &i
	return nil
}

type Maximum int

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
	i := int64(m)
	schema.Maximum = &i
	return nil
}

type ExclusiveMinimum int

func (e ExclusiveMinimum) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Exclusive minimum",
		Default: "",
	}
}

type ExclusiveMaximum int

func (e ExclusiveMaximum) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Exclusive maximum",
		Default: "",
	}
}
