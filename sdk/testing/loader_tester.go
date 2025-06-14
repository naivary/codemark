package testing

import (
	"reflect"
	"testing"
	"unicode"

	"github.com/naivary/codemark/parser"
)

type Struct struct {
	name    string
	markers []parser.Marker
	fields  []Field
}

type Field struct {
	f       reflect.StructField
	markers []parser.Marker
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
	return Struct{
		name:    randName(),
		fields:  fields,
		markers: randMarkers(),
	}
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
		f: reflect.StructField{
			Name: randName(),
			Type: randType(),
		},
		markers: randMarkers(),
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
