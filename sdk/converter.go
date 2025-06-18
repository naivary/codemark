package sdk

import (
	"reflect"

	"github.com/naivary/codemark/parser"
)

type TypeID = string

type Converter interface {
	// SupportedTypes returns the supported types for which the converter can be
	// used.
	SupportedTypes() []reflect.Type

	// CanConvert validates if the conversion of the marker to `def.output` is
	// possible. You can be sure that the convert is choosen correctly by the
	// ConverterManager and do not have to check if the `def.Output` is
	// convertible using this converter.
	CanConvert(m parser.Marker, def *Definition) error

	// Convert converts the marker to `def.Output`
	Convert(m parser.Marker, def *Definition) (reflect.Value, error)
}

type ConverterManager interface {
	GetConverter(rtype reflect.Type) (Converter, error)

	AddConverter(conv Converter) error

	AddConvByRefType(rtype reflect.Type, conv Converter) error

	Convert(mrk parser.Marker, target Target) (any, error)

	ParseDefs(doc string, t Target) (map[string][]any, error)
}
