package openapi

import (
	"go/types"
	"net/url"
)

type Schema struct {
	// metadata --
	ID    string   `json:"$id,omitzero"`
	Draft string   `json:"$schema,omitzero"`
	Type  jsonType `json:"type"`

	// annotations --
	Title      string   `json:"title,omitzero"`
	Desc       string   `json:"description,omitzero"`
	Examples   []string `json:"examples,omitzero"`
	Deprecated bool     `json:"deprecated,omitzero"`
	WriteOnly  bool     `json:"writeOnly,omitzero"`
	ReadOnly   bool     `json:"readOnly,omitzero"`
	Default    string   `json:"default,omitzero"`

	// object --
	Properties map[string]*Schema `json:"properties,omitzero"`
	Required   []string           `json:"required,omitzero"`

	// string --
	MinLength        int    `json:"minLength,omitzero"`
	MaxLength        int    `json:"maxLength,omitzero"`
	Pattern          string `json:"pattern,omitzero"`
	ContentEncoding  string `json:"contentEnconding,omitzero"`
	ContentMediaType string `json:"contentMediaType,omitzero"`
	Format           string `json:"format,omitzero"`

	// number --
	Maximum          int64 `json:"maximum,omitzero"`
	Minimum          int64 `json:"minimum,omitzero"`
	ExclusiveMaximum int64 `json:"exclusiveMaximum,omitzero"`
	ExclusiveMinimum int64 `json:"exclusiveMinimum,omitzero"`
	MultipleOf       int64 `json:"multipleOf,omitzero"`
}

type jsonType string

const (
	invalidType jsonType = "invalid"
	nullType    jsonType = "null"
	booleanType jsonType = "boolean"
	numberType  jsonType = "number"
	integerType jsonType = "integer"
	stringType  jsonType = "string"
	arrayType   jsonType = "array"
	objectType  jsonType = "object"
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

func id(name, baseURL string, nc NamingConvention) (string, error) {
	id := nc.Format(name) + ".json"
	idURL, err := url.JoinPath(baseURL, id)
	if err != nil {
		return "", err
	}
	return idURL, nil
}
