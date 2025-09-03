package openapi

import (
	"errors"
	"fmt"
	"regexp"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

type Pattern string

func (p Pattern) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Pattern which the given string must fullfill",
		Default: "",
		Type:    "string",
	}
}

func (p Pattern) apply(schema *Schema) error {
	regExp := string(p)
	if !p.isValidRegExp() {
		return fmt.Errorf("pattern is not a valid regular expression: %s", regExp)
	}
	if schema.Type == stringType {
		schema.Pattern = regExp
		return nil
	}
	if schema.Items.Type == stringType {
		schema.Items.Pattern = regExp
		return nil
	}
	if schema.AdditionalProperties.Type == stringType {
		schema.AdditionalProperties.Pattern = regExp
		return nil
	}
	return errors.New("pattern is only appliable to strings, objects with value of type string and arrays/slices of type string")
}

func (p Pattern) isValidRegExp() bool {
	_, err := regexp.Compile(string(p))
	return err == nil
}

type MinLength int64

func (m MinLength) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Minimum length of the string",
		Default: "",
	}
}

func (m MinLength) apply(schema *Schema) error {
	minLen := int64(m)
	if minLen < 0 {
		return fmt.Errorf("minLength cannot be negative: %d", minLen)
	}
	if schema.Type == stringType {
		schema.MinLength = minLen
		return nil
	}
	if schema.Items.Type == stringType {
		schema.Items.MinLength = minLen
		return nil
	}
	if schema.AdditionalProperties.Type == stringType {
		schema.AdditionalProperties.MinLength = minLen
		return nil
	}
	return errors.New("minLength is only appliable to strings, objects with value of type string and arrays/slices of type string")
}

type MaxLength int

func (m MaxLength) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Maximum length of the string",
		Default: "",
	}
}

func (m MaxLength) apply(schema *Schema) error {
	maxLen := int64(m)
	if maxLen < 0 {
		return fmt.Errorf("minLength cannot be negative: %d", maxLen)
	}
	if schema.Type == stringType {
		schema.MaxLength = maxLen
		return nil
	}
	if schema.Items.Type == stringType {
		schema.Items.MaxLength = maxLen
		return nil
	}
	if schema.AdditionalProperties.Type == stringType {
		schema.AdditionalProperties.MaxLength = maxLen
		return nil
	}
	return errors.New("maxLength is only appliable to strings, objects with value of type string and arrays/slices of type string")
}

type ContentEncoding string

func (c ContentEncoding) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Content encoding",
		Default: "",
	}
}

func (c ContentEncoding) apply(schema *Schema) error {
	encoding := string(c)
	if len(c) == 0 {
		return errors.New("contentEncoding marker cannot be empty")
	}
	if schema.Type == stringType {
		schema.ContentEncoding = encoding
		return nil
	}
	if schema.Items.Type == stringType {
		schema.Items.ContentEncoding = encoding
		return nil
	}
	if schema.AdditionalProperties.Type == stringType {
		schema.AdditionalProperties.ContentEncoding = encoding
		return nil
	}
	return errors.New("contentEncoding is only appliable to strings, objects with value of type string and arrays/slices of type string")
}

type ContentMediaType string

func (c ContentMediaType) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Content media type",
		Default: "",
	}
}

func (c ContentMediaType) apply(schema *Schema) error {
	mediaType := string(c)
	if len(mediaType) == 0 {
		return errors.New("contentMediaType cannot be empty")
	}
	if schema.Type == stringType {
		schema.ContentMediaType = mediaType
		return nil
	}
	if schema.Items.Type == stringType {
		schema.Items.ContentMediaType = mediaType
		return nil
	}
	if schema.AdditionalProperties.Type == stringType {
		schema.AdditionalProperties.ContentMediaType = mediaType
		return nil
	}
	return errors.New("contentMediaType is only appliable to strings, objects with value of type string and arrays/slices of type string")
}

type Format string

func (f Format) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Format",
		Default: "",
	}
}

func (f Format) apply(schema *Schema) error {
	format := string(f)
	if len(format) == 0 {
		return errors.New("format marker cannot be empty")
	}
	if schema.Type == stringType {
		schema.Format = format
		return nil
	}
	if schema.Items.Type == stringType {
		schema.Items.Format = format
		return nil
	}
	if schema.AdditionalProperties.Type == stringType {
		schema.AdditionalProperties.Format = format
		return nil
	}
	return errors.New("format is only appliable to strings, objects with value of type string and arrays/slices of type string")
}
