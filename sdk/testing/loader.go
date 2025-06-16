package testing

import (
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"
	"unicode"

	"github.com/naivary/codemark/parser"
)

type LoaderTestCase struct {
	Dir     string
	Structs []Struct
	Funcs   []Func
}

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

func NewGoFiles() (LoaderTestCase, error) {
	tc := LoaderTestCase{}
	structQuantity := quantity(6)
	funcQuantity := quantity(4)
	for range structQuantity {
		tc.Structs = append(tc.Structs, RandStruct())
	}
	for range funcQuantity {
		tc.Funcs = append(tc.Funcs, randFunc())
	}
	tmpl, err := template.ParseGlob("sdk/testing/tmpl/*")
	if err != nil {
		return tc, err
	}
	tmpDir := os.TempDir()
	dir, err := os.MkdirTemp(tmpDir, "cm-project")
	if err != nil {
		return tc, err
	}
	tc.Dir = dir
	for _, t := range tmpl.Templates() {
		name, _ := strings.CutSuffix(t.Name(), ".tmpl")
		path := path.Join(dir, name)
		file, err := os.Create(path)
		if err != nil {
			return tc, err
		}
		if err := t.Execute(file, &tc); err != nil {
			return tc, err
		}
	}
	return tc, err
}

func RandStruct() Struct {
	fieldQuantity := (randInt64() % 6) + 1
	methodQuantity := (randInt64() % 2) + 1
	fields := make([]Field, 0, fieldQuantity)
	methods := make([]Func, 0, methodQuantity)
	for range fieldQuantity {
		fields = append(fields, randField())
	}
	for range methodQuantity {
		methods = append(methods, randFunc())
	}
	s := Struct{
		Name:    randName(),
		Fields:  fields,
		Markers: randMarkers(),
		Methods: methods,
	}
	return s
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

func randFunc() Func {
	fn := reflect.FuncOf([]reflect.Type{}, []reflect.Type{}, false)
	return Func{
		Name:    randName(),
		Fn:      fn,
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

func quantity(mx int) int {
	q := (randInt64() % int64(mx)) + 1
	return int(q)
}
