package sdk

import (
	"reflect"

	"github.com/naivary/codemark/parser"
)

type ConverterManager interface {
	Convert(mrk parser.Marker, target Target) (any, error)

	AddConverter(conv Converter) error

	GetConverter(rtype reflect.Type) (Converter, error)
}
