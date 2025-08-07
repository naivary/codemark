package openapi

import (
	"errors"
	"go/types"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

type Enum []any

func (e Enum) Doc() docv1.Option {
	return docv1.Option{}
}

func (e Enum) apply(schema *Schema, obj types.Object) error {
	typ := e.typeOf(obj.Type())
	if typ == nil {
		return errors.New(
			"if this error message appears something is wrong with the decition making intenrally. Open an issue showing the request and your struct with the given marker",
		)
	}
	if _, isIface := typ.(*types.Interface); isIface {
		e.assign(schema)
		return nil
	}
	basic := typ.(*types.Basic)
	var err error
	switch basic.Kind() {
	case types.String:
		err = isTypeT[string](e)
	case types.Int, types.Int16, types.Int32, types.Int64, types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
		err = isTypeT[int64](e)
	case types.Float32, types.Float64:
		err = isTypeT[float64](e)
	case types.Bool:
		err = isTypeT[bool](e)
	}
	e.assign(schema)
	return err
}

func (e Enum) assign(schema *Schema) {
	for i, elem := range e {
		if elem == "null" {
			e[i] = nil
		}
	}
	if schema.Type == arrayType {
		schema.Items.Enum = e
		return
	}
	schema.Enum = e
}

func (e Enum) typeOf(typ types.Type) types.Type {
	iface, isIface := typ.(*types.Interface)
	if isIface && iface.Empty() {
		return iface
	}
	alias, isAlias := typ.(*types.Alias)
	if isAlias {
		return e.typeOf(alias.Rhs())
	}
	basic, isBasic := typ.(*types.Basic)
	if isBasic {
		return basic
	}
	slice, isSlice := typ.(*types.Slice)
	if isSlice {
		return e.typeOf(slice.Elem())
	}
	array, isArray := typ.(*types.Array)
	if isArray {
		return e.typeOf(array.Elem())
	}
	return nil
}

func isTypeT[T any](s []any) error {
	for _, el := range s {
		_, isT := el.(T)
		if !isT && el != "null" {
			return errors.New("not the same type")
		}
	}
	return nil
}
