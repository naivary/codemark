package openapi

import (
	"errors"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

type Pattern string

func (p Pattern) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Pattern which the given string must fullfill",
		Default: "",
	}
}

func (p Pattern) apply(schema *Schema) error {
	if schema.Type != stringType {
		return errors.New("pattern can only be applied to string types")
	}
	schema.Pattern = string(p)
	return nil
}

type MinLength int

func (m MinLength) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Minimum length of the string",
		Default: "",
	}
}

func (m MinLength) apply(schema *Schema) error {
	if schema.Type != stringType {
		return errors.New("minLength can only be applied to string types")
	}
	schema.MinLength = int(m)
	return nil
}

type MaxLength int

func (m MaxLength) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Maximum length of the string",
		Default: "",
	}
}

func (m MaxLength) apply(schema *Schema) error {
	if schema.Type != stringType {
		return errors.New("maxLength can only be applied to string types")
	}
	schema.MaxLength = int(m)
	return nil
}

type ContentEncoding string

func (c ContentEncoding) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Content encoding",
		Default: "",
	}
}

func (c ContentEncoding) apply(schema *Schema) error {
	if schema.Type != stringType {
		return errors.New("contentEncoding can only be applied to string types")
	}
	schema.ContentEncoding = string(c)
	return nil
}

type ContentMediaType string

func (c ContentMediaType) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Content media type",
		Default: "",
	}
}

func (c ContentMediaType) apply(schema *Schema) error {
	if schema.Type != stringType {
		return errors.New("contentMediaType can only be applied to string types")
	}
	schema.ContentMediaType = string(c)
	return nil
}

type Format string

func (f Format) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Format",
		Default: "",
	}
}

func (f Format) apply(schema *Schema) error {
	if schema.Type == stringType {
		schema.Format = string(f)
		return nil
	}
	if schema.Type == arrayType && schema.Items.Type == stringType {
		schema.Items.Format = string(f)
		return nil
	}
	return errors.New("format can only be applied to string types and string arrays/slices")
}
