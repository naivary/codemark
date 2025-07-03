package testing

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/maker"
	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk"
	"github.com/naivary/codemark/sdk/utils"
)

// TODO: From -> To is kinda weird because it can only be from kind of marker
// value type to a custom type

// ValidValueFunc defines the logic to check if the converter correctly
// converted the value of the marker to the custom type.
type ValidValueFunc func(got, want reflect.Value) bool

type ConverterTestCase struct {
	// Name of the test case
	Name string
	// Marker to convert by the converter
	Marker marker.Marker
	// Target of the marker
	Target target.Target
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
	NewTest(from, to reflect.Type, isValidCase bool, t target.Target) (ConverterTestCase, error)

	MustNewTest(from, to reflect.Type, isValidCase bool, t target.Target) ConverterTestCase

	// AddVVFunc defines a ValidValueFunc for an example type to which a
	// supported type of the converter will be converted.
	AddVVFunc(to reflect.Type, fn ValidValueFunc) error

	// ValidTests returns all test cases which must be valid based on the
	// supported types of the converter.
	ValidTests() ([]ConverterTestCase, error)

	// Run runs the given test case against and checks if the conversion was
	// successful.
	Run(t *testing.T, tc ConverterTestCase)
}

var _ ConverterTester = (*converterTester)(nil)

type converterTester struct {
	conv  sdk.Converter
	vvfns map[reflect.Type]ValidValueFunc
}

// NewConvTester returns a new ConverterTester for the given converter. The
// parameter `toTypes` is providing a map of a supported type to an example
// custom type which is being used by the converter as a test. For example the
// built in integer converter is converter an int to a type Int int.
func NewConvTester(conv sdk.Converter) (ConverterTester, error) {
	c := &converterTester{
		conv:  conv,
		vvfns: make(map[reflect.Type]ValidValueFunc),
	}
	return c, nil
}

// from should only be customizable for list others are
func (c *converterTester) NewTest(from, to reflect.Type, isValidCase bool, t target.Target) (ConverterTestCase, error) {
	marker, err := RandMarker(from)
	if err != nil {
		return ConverterTestCase{}, err
	}
	if to == nil {
		return ConverterTestCase{}, errors.New("to cannot be nil")
	}
	name := fmt.Sprintf("marker[%s] to %v", marker.Ident, to)
	vvfn := c.vvfns[to]
	if vvfn == nil {
		return ConverterTestCase{}, fmt.Errorf("ValidValueFunc not found: %v\n", from)
	}
	tc := ConverterTestCase{
		Name:         name,
		Marker:       *marker,
		Target:       t,
		To:           to,
		IsValidCase:  isValidCase,
		IsValidValue: vvfn,
	}
	return tc, nil
}

func (c *converterTester) MustNewTest(from, to reflect.Type, isValidCase bool, t target.Target) ConverterTestCase {
	tc, err := c.NewTest(from, to, isValidCase, t)
	if err != nil {
		panic(err)
	}
	return tc
}

func (c *converterTester) Run(t *testing.T, tc ConverterTestCase) {
	fakeDef, err := maker.MakeFakeDef(tc.To)
	if err != nil {
		t.Fatalf("err occured: %s\n", err)
	}
	v, err := c.conv.Convert(tc.Marker, fakeDef)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	gotType := reflect.TypeOf(v)
	if gotType != tc.To {
		t.Fatalf("types don't match after conversion. got: %v; want: %v\n", gotType, tc.To)
	}
	gotValue := reflect.ValueOf(v)
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

func (c *converterTester) ValidTests() ([]ConverterTestCase, error) {
	types := c.conv.SupportedTypes()
	tests := make([]ConverterTestCase, 0, len(types))
	for _, to := range types {
		mkind := utils.MarkerKindOf(to)
		tc, err := c.NewTest(mkind, to, true, target.ANY)
		if err != nil {
			return nil, err
		}
		tests = append(tests, tc)
	}
	return tests, nil
}
