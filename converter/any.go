package converter

import (
	"reflect"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	"github.com/naivary/codemark/marker"
)

var _ convv1.Converter = (*anyConverter)(nil)

type anyConverter struct {
	name string
}

func NewAny() convv1.Converter {
	return &anyConverter{
		name: "any",
	}
}

func (a *anyConverter) Name() string {
	return NewName(_codemark, a.name)
}

func (a *anyConverter) SupportedTypes() []reflect.Type {
	supported := []reflect.Type{
		reflect.TypeFor[any](),
		reflect.TypeFor[*any](),
	}
	return supported
}

func (a *anyConverter) CanConvert(m marker.Marker, to reflect.Type) error {
	return nil
}

func (a *anyConverter) Convert(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	return ConvertTo(m.Value, to)
}
