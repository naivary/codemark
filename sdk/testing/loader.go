package testing

import (
	"io"
	"reflect"

	"github.com/naivary/codemark/parser"
)

type Func struct {
	Name    string
	Fn      reflect.Type
	Markers []parser.Marker
}

type Struct struct {
	Name    string
	Markers []parser.Marker
	Fields  []Field
	Methods []Func
}

type Field struct {
	F       reflect.StructField
	Markers []parser.Marker
}

type LoaderTestCase struct {
	Structs []Struct
	Funcs   []Func
}

type LoaderTester interface {
	NewFile() (io.Reader, error)
}
