package testing

import (
	"reflect"

	"github.com/naivary/codemark/parser"
	"github.com/spf13/afero"
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
	NewFS() (afero.Fs, error)
}
