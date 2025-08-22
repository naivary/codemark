package openapi

import (
	"errors"
	"go/types"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/typeutil"
)

type Enum []any

func (e Enum) Doc() docv1.Option {
	return docv1.Option{}
}

func (e Enum) apply(schema *Schema, obj types.Object) error {
	if len(e) == 0 {
		return errors.New("enum cannot be empty")
	}
	err := marker.IsTypedList(obj.Type(), e)
	if err != nil {
		return err
	}
	basic := typeutil.BasicTypeOf(obj.Type())
	switch basic.Kind() {
	case types.Bool:
		return errors.New("an enum for a boolean type(primitive, array, slice, map etc.) is unnecessary")
	}
	e.assign(schema)
	return err
}

func (e Enum) assign(schema *Schema) {
	// this for loop is allowing null to be set in any kind of array and is
	// mapping it to the native null type of JSON. This is because codemark does
	// not know the null keyword.
	for i, elem := range e {
		if elem == "nil" {
			e[i] = nil
		}
	}
	if schema.Type == arrayType {
		schema.Items.Enum = e
		return
	}
	schema.Enum = e
}
