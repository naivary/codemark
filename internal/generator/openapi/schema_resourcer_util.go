package openapi

import (
	"go/types"
)

type Schema struct {
	ID    string   `json:"$id"`
	Draft string   `json:"$schema"`
	Desc  string   `json:"description"`
	Type  jsonType `json:"type"`

	// object --
	Properties map[string]*Schema

	// string --
	MinLength        *int
	MaxLength        *int
	Pattern          *string
	ContentEncoding  *string
	ContentMediaType *string

	// number --
	Maximum          *int64
	Minimum          *int64
	ExclusiveMaximum *int64
	ExclusiveMinimum *int64
	MultipleOf       *int64
}

type jsonType int

const (
	invalidType jsonType = iota
	nullType
	booleanType
	numberType
	integerType
	stringType
	arrayType
	objectType
)

func jsonTypeOf(typ types.Type) jsonType {
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return booleanType
		case types.String:
			return stringType
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			return integerType
		case types.Float32, types.Float64:
			return numberType
		case types.UntypedNil:
			return nullType
		default:
			return invalidType
		}

	case *types.Slice, *types.Array:
		return arrayType

	case *types.Map:
		// JSON only allows maps with string keys
		if keyType := t.Key(); keyType.Underlying().(*types.Basic).Kind() == types.String {
			return objectType
		}
		return invalidType

	case *types.Struct:
		return objectType

	case *types.Interface:
		// Accept interface{} as a valid JSON-compatible value
		return nullType

	case *types.Pointer:
		// Dereference and check underlying type
		return jsonTypeOf(t.Elem())
	default:
		return invalidType
	}
}

func (j jsonType) String() string {
	switch j {
	case nullType:
		return "null"
	case booleanType:
		return "boolean"
	case numberType:
		return "number"
	case integerType:
		return "integer"
	case stringType:
		return "string"
	case arrayType:
		return "array"
	case objectType:
		return "object"
	default:
		return ""
	}
}
