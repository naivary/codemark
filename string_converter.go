package codemark

import "github.com/naivary/codemark/marker"

var _ Converter = (*stringConverter)(nil)

type stringConverter struct {
	reg Registry
}

func (s *stringConverter) Convert(marker marker.Marker, target Target) (any, error) {
	return nil, nil
}
