package sdk

import (
	"reflect"

	"github.com/naivary/codemark/definition"
	"github.com/naivary/codemark/parser/marker"
)

type Converter interface {
	// Name is a human friendly representation of the converter. It may or may
	// not be unique depending on the Converter. It is recommended to prefix the
	// name of the converter with your project. For example the builtin
	// converters are prefixed with `codemark.*`.
	Name() string

	// SupportedTypes returns the types to which the Converter can converter
	// given a correct marker. If a marker is convertible to a supported type
	// can be validated using `CanConvert`.
	SupportedTypes() []reflect.Type

	// CanConvert validates if the conversion of the marker to `def.output` is
	// possible. You can be sure that the convert is choosen correctly by the
	// ConverterManager and do not have to check if the `def.Output` is
	// convertible using this converter.
	CanConvert(m marker.Marker, def *definition.Definition) error

	// Convert converts the marker to `def.Output`
	Convert(m marker.Marker, def *definition.Definition) (reflect.Value, error)
}
