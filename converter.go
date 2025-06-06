package codemark

import "github.com/naivary/codemark/parser"

type Converter interface {
	SupportedTypes() []any

	// CanConvert validates if the conversion of the marker to `def.output` is
	// possible.
	CanConvert(m parser.Marker, def *Definition) error

	// Convert converts the marker to `def.output`
	Convert(m parser.Marker, def *Definition) (any, error)
}
