package converter

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/marker"
)

var _ converter.Converter = (*boolConverter)(nil)

type boolConverter struct{}

func Bool() converter.Converter {
	return &boolConverter{}
}

func (b *boolConverter) Name() string {
	return converter.NewName(_codemark, "bool")
}

func (b *boolConverter) SupportedTypes() []reflect.Type {
	types := []any{
		bool(false),
		// pointer
		new(bool),
	}
	supported := make([]reflect.Type, 0, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		supported = append(supported, rtype)
	}
	return supported
}

func (b *boolConverter) CanConvert(m marker.Marker, to reflect.Type) error {
	if m.Kind != marker.BOOL {
		return fmt.Errorf(
			"marker kind of `%s` cannot be converted to a boolean. valid option is: %s",
			m.Kind,
			marker.BOOL,
		)
	}
	return nil
}

func (b *boolConverter) Convert(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	return converter.ConvertTo(m.Value, to)
}
