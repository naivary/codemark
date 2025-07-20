package converter

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/marker"
)

var _ converter.Converter = (*stringConverter)(nil)

type stringConverter struct {
	name string
}

func String() converter.Converter {
	return &stringConverter{
		name: "string",
	}
}

func (s *stringConverter) Name() string {
	return converter.NewName(_codemark, s.name)
}

func (s *stringConverter) SupportedTypes() []reflect.Type {
	types := []any{
		string(""),
		// pointer
		new(string),
	}
	supported := make([]reflect.Type, 0, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		supported = append(supported, rtype)
	}
	return supported
}

func (s *stringConverter) CanConvert(m marker.Marker, to reflect.Type) error {
	if m.Kind != marker.STRING {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a string. valid option is: %s", m.Kind, marker.STRING)
	}
	return nil
}

func (s *stringConverter) Convert(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	return converter.ConvertTo(m.Value, to)
}
