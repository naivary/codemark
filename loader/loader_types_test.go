package loader

import (
	"reflect"

	"github.com/naivary/codemark/parser/marker"
)

type markers interface {
	markers() []marker.Marker
}

type File struct {
	Markers []marker.Marker
}

func (f File) markers() []marker.Marker {
	return f.Markers
}

type Import struct {
	Name    string
	Markers []marker.Marker
}

func (i Import) markers() []marker.Marker {
	return i.Markers
}

type Alias struct {
	Name    string
	Type    reflect.Type
	Markers []marker.Marker
}

func (a Alias) markers() []marker.Marker {
	return a.Markers
}

type Named struct {
	Name    string
	Type    reflect.Type
	Markers []marker.Marker

	Methods map[string]Func
}

func (n Named) markers() []marker.Marker {
	return n.Markers
}

type Func struct {
	Name    string
	Fn      reflect.Type
	Markers []marker.Marker
}

func (f Func) markers() []marker.Marker {
	return f.Markers
}

type Struct struct {
	Name    string
	Markers []marker.Marker
	Fields  map[string]Field
	Methods map[string]Func
}

func (s Struct) markers() []marker.Marker {
	return s.Markers
}

type Field struct {
	F       reflect.StructField
	Markers []marker.Marker
}

func (f Field) markers() []marker.Marker {
	return f.Markers
}

type Const struct {
	Name    string
	Value   int64
	Markers []marker.Marker
}

func (c Const) markers() []marker.Marker {
	return c.Markers
}

type Var struct {
	Name    string
	Value   int64
	Markers []marker.Marker
}

func (v Var) markers() []marker.Marker {
	return v.Markers
}

type Iface struct {
	Name       string
	Signatures map[string]Func
	Markers    []marker.Marker
}

func (i Iface) markers() []marker.Marker {
	return i.Markers
}
