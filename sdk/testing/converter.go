package testing

import (
	"reflect"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
)

type ConverterTestCase struct {
	Name         string
	Marker       parser.Marker
	Target       sdk.Target
	ToType       reflect.Type
	IsValidCase  bool
	IsValidValue func(got reflect.Value, wanted reflect.Value) bool
}

func NewConvTestCase(name string, target sdk.Target) {}
