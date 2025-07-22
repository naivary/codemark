package convertertest

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/naivary/codemark/api/core"
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
	// Target of the marker
	Target core.Target
	// Type to convert the marker to.
	To reflect.Type
	// If the test case is a valid or invalid case
	IsValidCase bool
	// Function to validate the value of the converter (after conversion) with
	// the value of the given marker.
	IsEqual func(got, want reflect.Value) bool
}

// Tester is providing useful abstractions for testing a converter in
// a convenient and easy way.
type Tester interface {
	// NewCase returns a new test case.
	NewCase(to reflect.Type, isValidCase bool, t core.Target) (Case, error)

	// NewCaseWithMarker returns a test case with a custom marker.
	NewCaseWithMarker(
		marker *marker.Marker,
		to reflect.Type,
		isValidCase bool,
		t core.Target,
	) (Case, error)

	MustNewCase(to reflect.Type, isValidCase bool, t core.Target) Case
	MustNewCaseWithMarker(
		marker *marker.Marker,
		to reflect.Type,
		isValidCase bool,
		t core.Target,
	) Case

	// AddVVFunc defines a func(got, want reflect.Value) bool for an example type to which a
	// supported type of the converter will be converted.
	AddEqualFunc(to reflect.Type, fn func(got, want reflect.Value) bool) error

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

// from should only be customizable for list others are
func (c *tester) NewCase(to reflect.Type, isValidCase bool, t core.Target) (Case, error) {
	marker, err := markertest.RandMarker(to)
	if err != nil {
		return _casez, err
	}
	return c.NewCaseWithMarker(marker, to, isValidCase, t)
}

func (c *tester) NewCaseWithMarker(
	m *marker.Marker,
	to reflect.Type,
	isValidCase bool,
	t core.Target,
) (Case, error) {
	if m == nil {
		return _casez, errors.New("marker cannot be nil. use NewTest if you want a random marker")
	}
	if to == nil {
		return _casez, errors.New("to cannot be nil")
	}
	name := fmt.Sprintf("marker[%s] to %v", m.Ident, to)
	equal := c.equalFuncs[to]
	if equal == nil {
		return _casez, fmt.Errorf("func(got, want reflect.Value) bool not found for: %v", to)
	}
	tc := Case{
		Name:        name,
		Marker:      *m,
		Target:      t,
		To:          to,
		IsValidCase: isValidCase,
		IsEqual:     equal,
	}
	return tc, nil
}

func (c *tester) MustNewCase(to reflect.Type, isValidCase bool, t core.Target) Case {
	tc, err := c.NewCase(to, isValidCase, t)
	if err != nil {
		panic(err)
	}
	return tc
}

func (c *tester) MustNewCaseWithMarker(
	m *marker.Marker,
	to reflect.Type,
	isValidCase bool,
	t core.Target,
) Case {
	tc, err := c.NewCaseWithMarker(m, to, isValidCase, t)
	if err != nil {
		panic(err)
	}
	return tc
}

func (c *tester) Run(t *testing.T, tc Case) {
	// TODO: this has to work
	// err := c.conv.CanConvert(tc.Marker, tc.To)
	// if err != nil {
	// 	t.Errorf("err occured: %s\n", err)
	// }
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

func (c *tester) AddEqualFunc(to reflect.Type, fn func(got, want reflect.Value) bool) error {
	_, found := c.equalFuncs[to]
	if found {
		return fmt.Errorf("equal function already exists: %s", to)
	}
	c.equalFuncs[to] = fn
	return nil
}
