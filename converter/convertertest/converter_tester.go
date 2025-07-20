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

// ValidValueFunc defines the logic to check if the converter correctly
// converted the value of the marker to the custom type.
type ValidValueFunc func(got, want reflect.Value) bool

type ConverterTestCase struct {
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
	IsValidValue ValidValueFunc
}

// ConverterTester is providing useful abstractions for testing a converter in
// a convenient and easy way.
type ConverterTester interface {
	// NewTest returns a new ConverterTestCase.
	NewTest(to reflect.Type, isValidCase bool, t core.Target) (ConverterTestCase, error)

	NewTestWithMarker(marker *marker.Marker, to reflect.Type, isValidCase bool, t core.Target) (ConverterTestCase, error)

	MustNewTest(to reflect.Type, isValidCase bool, t core.Target) ConverterTestCase
	MustNewTestWithMarker(marker *marker.Marker, to reflect.Type, isValidCase bool, t core.Target) ConverterTestCase

	// AddVVFunc defines a ValidValueFunc for an example type to which a
	// supported type of the converter will be converted.
	AddVVFunc(to reflect.Type, fn ValidValueFunc) error

	// Run runs the given test case against and checks if the conversion was
	// successful.
	Run(t *testing.T, tc ConverterTestCase)
}

var _ ConverterTester = (*converterTester)(nil)

type converterTester struct {
	conv  converter.Converter
	vvfns map[reflect.Type]ValidValueFunc
}

// NewConvTester returns a new ConverterTester for the given converter. The
// parameter `toTypes` is providing a map of a supported type to an example
// custom type which is being used by the converter as a test. For example the
// built in integer converter is converter an int to a type Int int.
func NewConvTester(conv converter.Converter) (ConverterTester, error) {
	c := &converterTester{
		conv:  conv,
		vvfns: make(map[reflect.Type]ValidValueFunc),
	}
	return c, nil
}

// from should only be customizable for list others are
func (c *converterTester) NewTest(to reflect.Type, isValidCase bool, t core.Target) (ConverterTestCase, error) {
	marker, err := markertest.RandMarker(to)
	if err != nil {
		return ConverterTestCase{}, err
	}
	return c.NewTestWithMarker(marker, to, isValidCase, t)
}

func (c *converterTester) NewTestWithMarker(m *marker.Marker, to reflect.Type, isValidCase bool, t core.Target) (ConverterTestCase, error) {
	if m == nil {
		return ConverterTestCase{}, errors.New("marker cannot be nil. use NewTest if you want a random marker")
	}
	if to == nil {
		return ConverterTestCase{}, errors.New("to cannot be nil")
	}
	name := fmt.Sprintf("marker[%s] to %v", m.Ident, to)
	vvfn := c.vvfns[to]
	if vvfn == nil {
		return ConverterTestCase{}, fmt.Errorf("ValidValueFunc not found for: %v\n", to)
	}
	tc := ConverterTestCase{
		Name:         name,
		Marker:       *m,
		Target:       t,
		To:           to,
		IsValidCase:  isValidCase,
		IsValidValue: vvfn,
	}
	return tc, nil

}

func (c *converterTester) MustNewTest(to reflect.Type, isValidCase bool, t core.Target) ConverterTestCase {
	tc, err := c.NewTest(to, isValidCase, t)
	if err != nil {
		panic(err)
	}
	return tc
}

func (c *converterTester) MustNewTestWithMarker(m *marker.Marker, to reflect.Type, isValidCase bool, t core.Target) ConverterTestCase {
	tc, err := c.NewTestWithMarker(m, to, isValidCase, t)
	if err != nil {
		panic(err)
	}
	return tc
}

func (c *converterTester) Run(t *testing.T, tc ConverterTestCase) {
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
	if !tc.IsValidValue(gotValue, tc.Marker.Value) {
		t.Fatalf("value is not correct. got: %v; want: %v\n", gotValue, tc.Marker.Value)
	}
	t.Logf("succesfully converted. got: %v; want: %v\n", gotType, tc.To)
}

func (c *converterTester) AddVVFunc(to reflect.Type, fn ValidValueFunc) error {
	_, found := c.vvfns[to]
	if found {
		return fmt.Errorf("ValidValueFunc already exists: %s\n", to)
	}
	c.vvfns[to] = fn
	return nil
}
