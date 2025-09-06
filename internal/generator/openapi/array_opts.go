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
	return docv1.Option{
		Desc: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum consectetur, magna ut aliquet varius, tortor felis maximus massa, id faucibus nisi erat vitae velit. Sed tortor quam, efficitur et volutpat mattis, sodales sed leo. Vivamus sagittis ligula at dui ultricies luctus eget sed arcu. Morbi finibus pretium tellus eu euismod. Maecenas accumsan in odio nec imperdiet. Sed sodales sollicitudin ipsum nec euismod. Nullam eu rutrum lorem, vel molestie sem. Fusce eu magna non lacus elementum fringilla.

Donec euismod scelerisque dui vel accumsan. Sed non est non dui rutrum sollicitudin. Vivamus at posuere erat. Vivamus euismod lobortis enim, eget bibendum eros luctus vel. Mauris ornare egestas blandit. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Donec id finibus nibh, eget imperdiet libero. Nam quis commodo arcu. Sed ullamcorper vehicula tristique. Cras leo lacus, ullamcorper quis faucibus vel, imperdiet ultrices neque. Vestibulum tincidunt nisl auctor blandit aliquet. Duis ornare, nisi nec convallis cursus, velit velit auctor sapien, in dignissim lectus ex a turpis. Morbi sit amet tempus risus, sagittis hendrerit lacus. Aenean vitae tincidunt lectus, nec mollis elit.`,
	}
}

func (m MinItems) apply(schema *Schema) error {
	minItems := int64(m)
	if minItems < 0 {
		return errors.New("minItems cannot be negative")
	}
	if schema.Type == arrayType {
		schema.MinItems = minItems
		return nil
	}
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.Type == arrayType {
		schema.AdditionalProperties.MinItems = minItems
		return nil
	}
	return errors.New("minItems marker is only appliable to array/list types and objects with the value of arrays/slices")
}

type MaxItems int64

func (m MaxItems) Doc() docv1.Option {
	return docv1.Option{}
}

func (m MaxItems) apply(schema *Schema) error {
	maxItems := int64(m)
	if maxItems < 0 {
		return errors.New("maxItems cannot be negative")
	}
	if schema.Type == arrayType {
		schema.MaxItems = maxItems
		return nil
	}
	if schema.AdditionalProperties.Type == arrayType {
		schema.AdditionalProperties.MaxItems = maxItems
		return nil
	}
	return errors.New("maxItems marker is only appliable to array/list types and objects with the value of arrays/slices")
}

type UniqueItems bool

func (u UniqueItems) Doc() docv1.Option {
	return docv1.Option{}
}

func (u UniqueItems) apply(schema *Schema) error {
	uniqueItems := bool(u)
	if !uniqueItems {
		return errors.New("by default uniqueItems will be false. Remove the marker so the amount of markers is not increased unnecessarily")
	}
	if schema.Type != arrayType {
		schema.UniqueItems = uniqueItems
		return nil
	}
	if schema.AdditionalProperties.Type == arrayType {
		schema.AdditionalProperties.UniqueItems = uniqueItems
		return nil
	}
	return errors.New("uniqueItems marker is only appliable to array/list types and objects with the value of arrays/slices")
}
