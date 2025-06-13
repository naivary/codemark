package testing

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
)

type ValidValueFunc func(got, want reflect.Value) bool

type ConverterTesterConfig struct {
	ValidValueFuncs map[sdk.TypeID]ValidValueFunc
	Types           map[sdk.TypeID]reflect.Type
}

type ConverterTestCase struct {
	Name         string
	Marker       parser.Marker
	Target       sdk.Target
	ToType       reflect.Type
	IsValidCase  bool
	IsValidValue ValidValueFunc
}

type ConverterTester interface {
	NewTestCase(rtype reflect.Type, isValidCase bool) (ConverterTestCase, error)

	ValidTestCases() ([]ConverterTestCase, error)

	Run(t *testing.T, tc ConverterTestCase, mngr sdk.ConverterManager)
}
