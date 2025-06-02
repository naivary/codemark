package codemark

import "github.com/naivary/codemark/parser"

type Converter interface {
	// SupportedTypeIDs returns which typeIDs are supported by this converter
	// and can be expected to succesfully be converted.
	SupportedTypeIDs() []string

	// IsPossible returns an error if the conversion of the `MarkerKind` is not
	// possible to `def.Output()`
	CanConvert(m parser.Marker, def *Definition) error

	Convert(m parser.Marker, def *Definition) (any, error)
}
