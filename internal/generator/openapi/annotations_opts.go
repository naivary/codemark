package openapi

import (
	"errors"
	"go/types"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	"github.com/naivary/codemark/marker"
)

type Description string

func (d Description) Doc() docv1.Option {
	return docv1.Option{
		Desc: "Description",
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
		Desc: "Title",
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

type Examples []any

func (e Examples) Doc() docv1.Option {
	return docv1.Option{
		Desc: "Examples",
	}
}

func (e Examples) apply(schema *Schema, fieldType types.Type) error {
	// Missing type checking of the field and examples e.g. []string can only
	// contain string examples not integer too
	if len(e) == 0 {
		return errors.New("examples marker cannot be emtpy")
	}
	err := marker.IsTypedList(fieldType, e)
	if err != nil {
		return err
	}
	if schema.Type == objectType || schema.Type == arrayType {
		return errors.New("examples annotation cannot be set for objects and arrays/slices")
	}
	schema.Examples = e
	return nil
}

type Deprecated bool

func (d Deprecated) Doc() docv1.Option {
	return docv1.Option{
		Desc: "Deprecated",
	}
}

func (d Deprecated) apply(schema *Schema) error {
	schema.Deprecated = bool(d)
	return nil
}

type ReadOnly bool

func (r ReadOnly) Doc() docv1.Option {
	return docv1.Option{
		Desc: "indicates that the field is read only",
	}
}

func (r ReadOnly) apply(schema *Schema) error {
	if !r {
		return errors.New("readOnly is by default already false. Don't set the marker unnecessarily")
	}
	schema.ReadOnly = bool(r)
	return nil
}

type WriteOnly bool

func (w WriteOnly) Doc() docv1.Option {
	return docv1.Option{
		Desc: "indicates that the field is write only",
	}
}

func (w WriteOnly) apply(schema *Schema) error {
	if !w {
		return errors.New("writeOnly is by default already false. Don't set the marker unnecessarily")
	}
	schema.WriteOnly = bool(w)
	return nil
}

type Default string

func (d Default) Doc() docv1.Option {
	return docv1.Option{
		Desc: "default value fot this field",
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
