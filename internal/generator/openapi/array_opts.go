package openapi

import (
	"errors"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

// The options are left out because of the following reasons:
// 1. prefixItems: if you want a prefixed order use structs
// 2. unevaluatedItems: because its go
// 3. additionalItems: deprecated
//
// The following options are still under considertion:
// 1. contains (especially useful with []any slices/arrays)
// 2. minContains
// 3. maxContains

type MinItems int64

func (m MinItems) Doc() docv1.Option {
	return docv1.Option{}
}

func (m MinItems) apply(schema *Schema) error {
	if schema.Type != arrayType {
		return errors.New("minItems marker is only appliable to array/list types")
	}
	schema.MinItems = int64(m)
	return nil
}

type MaxItems int64

func (m MaxItems) Doc() docv1.Option {
	return docv1.Option{}
}

func (m MaxItems) apply(schema *Schema) error {
	if schema.Type != arrayType {
		return errors.New("maxItems marker is only appliable to array/list types")
	}
	schema.MaxItems = int64(m)
	return nil
}

type UniqueItems bool

func (u UniqueItems) Doc() docv1.Option {
	return docv1.Option{}
}

func (u UniqueItems) apply(schema *Schema) error {
	if schema.Type != arrayType {
		return errors.New("maxItems marker is only appliable to array/list types")
	}
	schema.UniqueItems = bool(u)
	return nil
}
