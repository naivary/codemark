package loader

import (
	"reflect"

	"github.com/naivary/codemark/parser"
)

type defs interface {
	Definitions() map[string][]any
}

type markers interface {
	markers() []parser.Marker
}

type File struct {
	Markers []parser.Marker
}

func (f File) markers() []parser.Marker {
	return f.Markers
}

type Import struct {
	Name    string
	Markers []parser.Marker
}

func (i Import) markers() []parser.Marker {
	return i.Markers
}

type Alias struct {
	Name    string
	Type    reflect.Type
	Markers []parser.Marker
}

func (a Alias) markers() []parser.Marker {
	return a.Markers
}

type Named struct {
	Name    string
	Type    reflect.Type
	Markers []parser.Marker

	Methods map[string]Func
}

func (n Named) markers() []parser.Marker {
	return n.Markers
}

type Func struct {
	Name    string
	Fn      reflect.Type
	Markers []parser.Marker
}

func (f Func) markers() []parser.Marker {
	return f.Markers
}

type Struct struct {
	Name    string
	Markers []parser.Marker
	Fields  map[string]Field
	Methods map[string]Func
}

func (s Struct) markers() []parser.Marker {
	return s.Markers
}

type Field struct {
	F       reflect.StructField
	Markers []parser.Marker
}

func (f Field) markers() []parser.Marker {
	return f.Markers
}

type Const struct {
	Name    string
	Value   int64
	Markers []parser.Marker
}

func (c Const) markers() []parser.Marker {
	return c.Markers
}

type Var struct {
	Name    string
	Value   int64
	Markers []parser.Marker
}

func (v Var) markers() []parser.Marker {
	return v.Markers
}

type Iface struct {
	Name       string
	Signatures map[string]Func
	Markers    []parser.Marker
}

func (i Iface) markers() []parser.Marker {
	return i.Markers
}
