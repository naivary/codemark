package codemark

import (
	"errors"

	"github.com/naivary/codemark/parser"
)

var _ Converter = (*stringConverter)(nil)

type stringConverter struct{}

func (s *stringConverter) SupportedTypes() []any {
	return []any{
		string(""),
		rune(0),
		byte(0),
		[]byte{},
		// pointer
		new(string),
		new(rune),
		new(byte),
		[]*byte{},
	}
}

func (s *stringConverter) CanConvert(m parser.Marker, def *Definition) error {
	if m.Kind() != parser.MarkerKindString {
		return errors.New("only string markers can be converted by the string converter")
	}
	return nil
}

func (s *stringConverter) Convert(m parser.Marker, def *Definition) (any, error) {
	return nil, nil
}
