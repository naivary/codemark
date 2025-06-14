package testing

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"text/template"
	"unicode"

	"github.com/naivary/codemark/parser"
)

type File struct {
	Structs []Struct
}

type Struct struct {
	Name    string
	Markers []parser.Marker
	Fields  []Field
}

type Field struct {
	F       reflect.StructField
	Markers []parser.Marker
}

var _ LoaderTester = (*loaderTester)(nil)

type loaderTester struct{}

func NewLoaderTester() (LoaderTester, error) {
	return nil, nil
}

func (l *loaderTester) Run(t *testing.T, tc LoaderTestCase) {}

func randStruct() Struct {
	fieldQuantity := (randInt64() % 6) + 1
	fields := make([]Field, 0, fieldQuantity)
	for range fieldQuantity {
		fields = append(fields, randField())
	}
	s := Struct{
		Name:    randName(),
		Fields:  fields,
		Markers: randMarkers(),
	}
	t, err := template.ParseGlob("tmpl/*")
	fmt.Println(err)
	err = t.Execute(os.Stdout, File{[]Struct{s}})
	fmt.Println(err)
	return s
}

func randName() string {
	name := randString()
	for {
		firstLetter := rune(name[0])
		if !unicode.IsDigit(firstLetter) {
			break
		}
		name = randString()
	}
	return name
}

func randField() Field {
	return Field{
		F: reflect.StructField{
			Name: randName(),
			Type: randType(),
		},
		Markers: randMarkers(),
	}
}

func randMarkers() []parser.Marker {
	markerQuantity := (randInt64() % 5) + 1
	markers := make([]parser.Marker, 0, markerQuantity)
	for range markerQuantity {
		markers = append(markers, RandMarker(randType()))
	}
	return markers
}

func randType() reflect.Type {
	// string, int, float32, complex64, bool, uint
	i := (randInt("int64")() % 11) + 1
	switch i {
	case 1:
		return reflect.TypeFor[string]()
	case 2:
		return reflect.TypeFor[int]()
	case 3:
		return reflect.TypeFor[float32]()
	case 4:
		return reflect.TypeFor[complex64]()
	case 5:
		return reflect.TypeFor[bool]()
	case 6:
		return reflect.TypeFor[uint]()
	case 7:
		return reflect.TypeFor[[]string]()
	case 8:
		return reflect.TypeFor[[]int]()
	case 9:
		return reflect.TypeFor[[]float32]()
	case 10:
		return reflect.TypeFor[[]complex64]()
	case 11:
		return reflect.TypeFor[[]bool]()
	}
	return nil
}
