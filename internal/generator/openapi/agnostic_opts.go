package openapi

import (
	"errors"
	"fmt"
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
		return fmt.Errorf(
			"if this error message appears something is wrong with the decition making intenrally. Open an issue showing the request and your struct with the used markers. struct: `%s`",
			obj.Name(),
		)
	}
	if _, isIface := typ.(*types.Interface); isIface {
		e.assign(schema)
		return nil
	}
	basic := typ.(*types.Basic)
	var err error
	// NOTE: the types (string, int64, float64 and bool) are all the types we
	// need to check because codemark is only using these types not the others.
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
	// this for loop is allowing null to be set in any kind of array and is
	// mapping it to the native null type of JSON. This is because codemark does
	// not know the null keyword.
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
	basic, isBasic := typ.(*types.Basic)
	if isBasic {
		return basic
	}

	switch t := typ.(type) {
	case *types.Alias:
		return e.typeOf(t.Rhs())
	case *types.Slice:
		return e.typeOf(t.Elem())
	case *types.Array:
		return e.typeOf(t.Elem())
	case *types.Named:
		return e.typeOf(t.Underlying())
	default:
		return nil
	}
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
