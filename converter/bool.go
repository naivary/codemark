package converter

import (
	"fmt"
	"reflect"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	"github.com/naivary/codemark/marker"
)

var _ convv1.Converter = (*boolConverter)(nil)

type boolConverter struct{}

func NewBool() convv1.Converter {
	return &boolConverter{}
}

func (b *boolConverter) Name() string {
	return NewName(_codemark, "bool")
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
	return ConvertTo(m.Value, to)
}
