package converter

import (
	"fmt"
	"reflect"
	"time"

	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/typeutil"
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
		return fmt.Errorf("marker kind of `%s` cannot be converted to a string. valid option is: %s", m.Kind, marker.STRING)
	}
	return nil
}

func (s *stringConverter) Convert(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	if !s.isTime(to) {
		return converter.ConvertTo(m.Value, to)
	}
	return s.time(m, to)
}

func (s *stringConverter) isTime(to reflect.Type) bool {
	to = typeutil.Deref(to)
	return to.ConvertibleTo(reflect.TypeOf(time.Time{}))
}

func (s *stringConverter) time(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	t, err := time.Parse(time.RFC3339, m.Value.String())
	if err != nil {
		return _rvzero, err
	}
	v := reflect.ValueOf(t)
	return converter.ConvertTo(v, to)
}
