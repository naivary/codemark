package converter

import (
	"fmt"
	"reflect"
	"time"

	"github.com/naivary/codemark/marker"
)

var _ Converter = (*stringConverter)(nil)

type stringConverter struct {
	name string
}

func String() Converter {
	return &stringConverter{
		name: "string",
	}
}

func (s *stringConverter) Name() string {
	return NewName(_codemark, s.name)
}

func (s *stringConverter) SupportedTypes() []reflect.Type {
	types := []any{
		string(""),
		time.Time{},
		// pointer
		new(string),
		new(time.Time),
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
		return fmt.Errorf(
			"marker kind of `%s` cannot be converted to a string. valid option is: %s",
			m.Kind,
			marker.STRING,
		)
	}
	return nil
}

func (s *stringConverter) Convert(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	if isTypeT[time.Time](to) {
		return s.time(m, to)
	}
	return ConvertTo(m.Value, to)
}

func (s *stringConverter) time(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	t, err := time.Parse(time.RFC3339, m.Value.String())
	if err != nil {
		return _rvzero, err
	}
	v := reflect.ValueOf(t)
	return ConvertTo(v, to)
}
