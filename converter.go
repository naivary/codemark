package codemark

import (
	"reflect"

	"github.com/naivary/codemark/parser"
)

type Converter interface {
	// SupportedTypes returns the supported types for which the converter can be
	// used.
	SupportedTypes() []reflect.Type

	// CanConvert validates if the conversion of the marker to `def.output` is
	// possible.
	CanConvert(m parser.Marker, def *Definition) error

	// Convert converts the marker to `def.output`
	Convert(m parser.Marker, def *Definition) (reflect.Value, error)
}
