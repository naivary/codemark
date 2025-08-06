package openapi

import (
	"errors"

	docv1 "github.com/naivary/codemark/api/doc/v1"
)

type Description string

func (d Description) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Description",
		Default: "",
	}
}

func (d Description) apply(schema *Schema) error {
	desc := string(d)
	if len(desc) == 0 {
		return errors.New("description for schema cannot be empty")
	}
	schema.Desc = desc
	return nil
}

type Title string

func (t Title) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Title",
		Default: "",
	}
}

func (t Title) apply(schema *Schema) error {
	str := string(t)
	if len(str) == 0 {
		return errors.New("title marker cannot be emtpy")
	}
	schema.Title = str
	return nil
}

type Examples []string

func (e Examples) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Examples",
		Default: "",
	}
}

func (e Examples) apply(schema *Schema) error {
	if len(e) == 0 {
		return errors.New("examples marker cannot be emtpy")
	}
	if schema.Type == objectType {
		return errors.New("examples annotation cannot be set for objects")
	}
	schema.Examples = e
	return nil
}

type Deprecated bool

func (d Deprecated) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "Deprecated",
		Default: "",
	}
}

func (d Deprecated) apply(schema *Schema) error {
	schema.Deprecated = bool(d)
	return nil
}

type ReadOnly bool

func (r ReadOnly) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "indicates that the field is read only",
		Default: "",
	}
}

func (r ReadOnly) apply(schema *Schema) error {
	schema.ReadOnly = bool(r)
	return nil
}

type WriteOnly bool

func (w WriteOnly) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "indicates that the field is write only",
		Default: "",
	}
}

func (w WriteOnly) apply(schema *Schema) error {
	schema.WriteOnly = bool(w)
	return nil
}

type Default string

func (d Default) Doc() docv1.Option {
	return docv1.Option{
		Desc:    "default value fot this field",
		Default: "",
	}
}

func (d Default) apply(schema *Schema) error {
	str := string(d)
	if len(str) == 0 {
		return errors.New("default marker cannot be emtpy")
	}
	if schema.Type == objectType {
		return errors.New("default annotation cannot be set for objects")
	}
	schema.Default = str
	return nil
}
