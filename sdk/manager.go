package sdk

import (
	"reflect"

	"github.com/naivary/codemark/parser"
)

type ConverterManager interface {
	AddConverter(conv Converter) error

	Convert(mrk parser.Marker, target Target) (any, error)

	GetConverter(rtype reflect.Type) (Converter, error)
}
