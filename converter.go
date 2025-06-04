package codemark

import "github.com/naivary/codemark/parser"

type Converter interface {
	// SupportedTypeIDs returns which typeIDs are supported by this converter
	// and can be expected to succesfully be converted. These will be added to
	// the default converters already present.
	SupportedTypeIDs() []string

	// CanConvert validates if the conversion of the marker to `def.output` is
	// possible.
	CanConvert(m parser.Marker, def *Definition) error

	// Convert converts the marker to `def.output`
	Convert(m parser.Marker, def *Definition) (any, error)
}
