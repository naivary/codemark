package testing

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
)

// ValidValueFunc defines the logic to check if the converter correctly
// converted the value of the marker to the custom type.
type ValidValueFunc func(got, want reflect.Value) bool

type ConverterTestCase struct {
	// Name of the test case
	Name string
	// Marker to convert by the converter
	Marker parser.Marker
	// Target of the marker
	Target sdk.Target
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
	NewTest(from reflect.Type, isValidCase bool, t sdk.Target) (ConverterTestCase, error)

	// AddVVFunc defines a ValidValueFunc for an example type to which a
	// supported type of the converter will be converted.
	AddVVFunc(to reflect.Type, fn ValidValueFunc) error

	// ValidTests returns all test cases which must be valid based on the
	// supported types of the converter.
	ValidTests() ([]ConverterTestCase, error)

	// Run runs the given test case against and checks if the conversion was
	// successful.
	Run(t *testing.T, tc ConverterTestCase, mngr sdk.ConverterManager)
}
