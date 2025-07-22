package convertertest

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/marker/markertest"
)

// _casez is the zero value
var _casez = Case{}

type Case struct {
	// Name of the test case
	Name string
	// Marker to convert by the converter
	Marker marker.Marker
	// Type to convert the marker to.
	To reflect.Type
	// If the test case is a valid or invalid case
	IsValidCase bool
	// Function to validate the value of the converter (after conversion) with
	// the value of the given marker.
	IsEqual func(got, want reflect.Value) bool
}

func NewCase(
	m *marker.Marker,
	to reflect.Type,
	isValidCase bool,
	equal func(got, want reflect.Value) bool,
) (Case, error) {
	var err error
	if m == nil {
		m, err = markertest.RandMarker(to)
	}
	if err != nil {
		return _casez, err
	}
	if to == nil {
		return _casez, errors.New("to cannot be nil")
	}
	if equal == nil {
		return _casez, fmt.Errorf("func(got, want reflect.Value) bool cannot be nil: %s", m)
	}

	name := fmt.Sprintf("marker[%s] to %v", m.Ident, to)
	tc := Case{
		Name:        name,
		Marker:      *m,
		To:          to,
		IsValidCase: isValidCase,
		IsEqual:     equal,
	}
	return tc, nil
}

func MustNewCase(
	m *marker.Marker,
	to reflect.Type,
	isValidCase bool,
	equal func(got, want reflect.Value) bool,
) Case {
	cas, err := NewCase(m, to, isValidCase, equal)
	if err != nil {
		panic(err)
	}
	return cas
}

// Tester is providing useful abstractions for testing a converter in
// a convenient and easy way.
type Tester interface {
	// NewCase returns a new test case.
	NewCase(
		m *marker.Marker,
		to reflect.Type,
		isValidCase bool,
		equal func(got, want reflect.Value) bool,
	) (Case, error)

	MustNewCase(
		m *marker.Marker,
		to reflect.Type,
		isValidCase bool,
		equal func(got, want reflect.Value) bool,
	) Case

	// Run runs the given test case against and checks if the conversion was
	// successful.
	Run(t *testing.T, tc Case)
}

var _ Tester = (*tester)(nil)

type tester struct {
	conv       converter.Converter
	equalFuncs map[reflect.Type]func(got, want reflect.Value) bool
}

// NewTester returns a new ConverterTester for the given converter. The
// parameter `toTypes` is providing a map of a supported type to an example
// custom type which is being used by the converter as a test. For example the
// built in integer converter is converter an int to a type Int int.
func NewTester(conv converter.Converter) (Tester, error) {
	c := &tester{
		conv:       conv,
		equalFuncs: make(map[reflect.Type]func(got, want reflect.Value) bool),
	}
	return c, nil
}

func (c *tester) NewCase(
	m *marker.Marker,
	to reflect.Type,
	isValidCase bool,
	equal func(got, want reflect.Value) bool,
) (Case, error) {
	return NewCase(m, to, isValidCase, equal)
}

func (c *tester) MustNewCase(
	m *marker.Marker,
	to reflect.Type,
	isValidCase bool,
	equal func(got, want reflect.Value) bool,
) Case {
	return MustNewCase(m, to, isValidCase, equal)
}

func (c *tester) Run(t *testing.T, tc Case) {
	err := c.conv.CanConvert(tc.Marker, tc.To)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	v, err := c.conv.Convert(tc.Marker, tc.To)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	value := v.Interface()
	gotType := reflect.TypeOf(value)
	if gotType != tc.To {
		t.Fatalf("types don't match after conversion. got: %v; want: %v\n", gotType, tc.To)
	}
	gotValue := reflect.ValueOf(value)
	if !tc.IsEqual(gotValue, tc.Marker.Value) {
		t.Fatalf("value is not correct. got: %v; want: %v\n", gotValue, tc.Marker.Value)
	}
	t.Logf("succesfully converted. got: %v; want: %v\n", gotType, tc.To)
}
