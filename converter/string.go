package converter

import (
	"errors"

	"github.com/naivary/codemark"
	"github.com/naivary/codemark/marker"
)

var _ Converter = (*stringConv)(nil)

type stringConv struct{}

func ForString() {}

func convertString(m marker.Marker, def *Definition) (any, error) {
	if !isStringConvPossible(def.kind) {
		return nil, errors.New("string conversion not possible")
	}
	s := m.Value().String()
	value, err := convString(s, def)
	if err != nil {
		return nil, err
	}
	return toOutput(value, def)
}

func (s *stringConv) Convert(m marker.Marker, t codemark.Target) (any, error) {
	return nil, nil
}

func (s *stringConv) CanConvert(m marker.Marker, def *codemark.Definition) error {
	// +path:to:marker=string -> string|*string & []byte|[]*byte
	// +path:to:marker=s -> rune|*rune & byte|*byte & string|*string
	return nil
}
