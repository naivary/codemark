package loader

import (
	"reflect"

	"github.com/naivary/codemark/parser"
)

type Import struct {
	Name    string
	Markers []parser.Marker
}

type Alias struct {
	Name    string
	Type    reflect.Type
	Markers []parser.Marker
}

type Named struct {
	Name    string
	Type    reflect.Type
	Markers []parser.Marker
}

type Func struct {
	Name    string
	Fn      reflect.Type
	Markers []parser.Marker
}

type Struct struct {
	Name    string
	Markers []parser.Marker
	Fields  map[string]Field
	Methods map[string]Func
}

type Field struct {
	F       reflect.StructField
	Markers []parser.Marker
}

type Const struct {
	Name    string
	Value   int64
	Markers []parser.Marker
}

type Var struct {
	Name    string
	Value   int64
	Markers []parser.Marker
}

type Iface struct {
	Name       string
	Signatures map[string]Func
	Markers    []parser.Marker
}
