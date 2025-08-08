package openapi

import (
	"errors"
	"fmt"
	"go/types"
	"net/url"
)

var _schemaz = Schema{}

type Schema struct {
	// metadata
	ID    string   `json:"$id,omitzero"`
	Draft string   `json:"$schema,omitzero"`
	Ref   string   `json:"$ref,omitzero"`
	Type  jsonType `json:"type,omitzero"`

	OneOf []*Schema `json:"oneOf,omitzero"`
	AnyOf []*Schema `json:"anyOf,omitzero"`
	Not   *Schema   `json:"not,omitzero"`

	// agnostic
	Enum []any `json:"enum,omitzero"`

	// annotations
	Title      string `json:"title,omitzero"`
	Desc       string `json:"description,omitzero"`
	Examples   []any  `json:"examples,omitzero"`
	Deprecated bool   `json:"deprecated,omitzero"`
	WriteOnly  bool   `json:"writeOnly,omitzero"`
	ReadOnly   bool   `json:"readOnly,omitzero"`
	Default    string `json:"default,omitzero"`

	// array
	MaxItems    int64   `json:"maxItems,omitzero"`
	MinItems    int64   `json:"minItems,omitzero"`
	UniqueItems bool    `json:"uniqueItems,omitzero"`
	Items       *Schema `json:"items,omitzero"`

	// object
	Properties           map[string]*Schema  `json:"properties,omitzero"`
	Required             []string            `json:"required,omitzero"`
	AdditionalProperties *Schema             `json:"additionalProperties,omitzero"`
	PatternProperties    map[string]*Schema  `json:"patternProperties,omitzero,omitempty"`
	DependentRequired    map[string][]string `json:"dependentRequired,omitzero,omitempty"`

	// string
	MinLength        int64  `json:"minLength,omitzero"`
	MaxLength        int64  `json:"maxLength,omitzero"`
	Pattern          string `json:"pattern,omitzero"`
	ContentEncoding  string `json:"contentEnconding,omitzero"`
	ContentMediaType string `json:"contentMediaType,omitzero"`
	Format           string `json:"format,omitzero"`

	// number
	Maximum          int64 `json:"maximum,omitzero"`
	Minimum          int64 `json:"minimum,omitzero"`
	ExclusiveMaximum int64 `json:"exclusiveMaximum,omitzero"`
	ExclusiveMinimum int64 `json:"exclusiveMinimum,omitzero"`
	MultipleOf       int64 `json:"multipleOf,omitzero"`
}

func newSchema(typ types.Type, cfg *config) (Schema, error) {
	// TODO: some markers make sense to be abel to set even if ref is active.
	// find those markers. For example maxItems, minItems.
	switch t := typ.(type) {
	case *types.Named:
		return newSchemaFromNamed(t, cfg)
	case *types.Alias:
		return newSchema(t.Rhs(), cfg)
	}

	switch t := typ.Underlying().(type) {
	case *types.Basic:
		return newBasicSchema(t)
	case *types.Slice:
		return newArraySchemaFromSlice(t, cfg)
	case *types.Array:
		return newArraySchemaFromArray(t, cfg)
	case *types.Map:
		return newObjectSchemaFromMap(t, cfg)
	case *types.Interface:
		return newSchemaFromIface(t)
	case *types.Pointer:
		return newSchema(t.Elem(), cfg)
	}
	return Schema{}, fmt.Errorf("type is not supported: %s", typ)
}

func newBasicSchema(t *types.Basic) (Schema, error) {
	return Schema{
		Type: jsonTypeOf(t),
	}, nil
}

func newSchemaFromIface(t *types.Interface) (Schema, error) {
	if !t.Empty() {
		return _schemaz, errors.New("interfaces are not supported as types in struct for JSON Schema beside the expection of any")
	}
	return Schema{}, nil
}

func newArraySchemaFromSlice(t *types.Slice, cfg *config) (Schema, error) {
	elemSchema, err := newSchema(t.Elem(), cfg)
	if err != nil {
		return _schemaz, err
	}
	schema := Schema{
		Type:  arrayType,
		Items: &elemSchema,
	}
	return schema, nil
}

func newArraySchemaFromArray(t *types.Array, cfg *config) (Schema, error) {
	elemSchema, err := newSchema(t.Elem(), cfg)
	if err != nil {
		return _schemaz, err
	}
	schema := Schema{
		Type:     arrayType,
		MaxItems: max(t.Len(), 0),
		Items:    &elemSchema,
	}
	return schema, nil
}

func newObjectSchemaFromMap(t *types.Map, cfg *config) (Schema, error) {
	basic, isBasic := t.Key().Underlying().(*types.Basic)
	if !isBasic {
		return _schemaz, errors.New("map is not indexed with an basic type e.g. int, string etc.")
	}
	if basic.Kind() != types.String {
		return _schemaz, errors.New("map has to be indexed with string")
	}
	valueSchema, err := newSchema(t.Elem(), cfg)
	if err != nil {
		return _schemaz, err
	}
	return Schema{
		Type:                 objectType,
		AdditionalProperties: &valueSchema,
	}, nil
}

func newObjectSchemaFromStruct(name string, cfg *config) (Schema, error) {
	id, err := id(name, cfg.Schema.IDBaseURL, cfg.Schema.Formats.Filename)
	if err != nil {
		return _schemaz, err
	}
	return Schema{
		Ref: id,
	}, nil
}

func newSchemaFromNamed(n *types.Named, cfg *config) (Schema, error) {
	switch n.Underlying().(type) {
	case *types.Struct:
		return newObjectSchemaFromStruct(n.Obj().Name(), cfg)
	}
	return _schemaz, nil
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
	case *types.Pointer:
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
