package converter

import (
	"reflect"

	"github.com/naivary/codemark/marker"
)

type Converter interface {
	// Name is a human friendly representation of the converter. It may or may
	// not be unique depending on the Converter. It is recommended to prefix the
	// name of the converter with your project. For example the builtin
	// converters are prefixed with `codemark.*`.
	Name() string

	// SupportedTypes returns the types to which the Converter can convert
	// given a correct marker.
	SupportedTypes() []reflect.Type

	// CanConvert validates if the conversion of the marker to the given type is
	// possible.
	CanConvert(m marker.Marker, to reflect.Type) error

	// Convert converts the marker to the given type
	Convert(m marker.Marker, to reflect.Type) (reflect.Value, error)
}
