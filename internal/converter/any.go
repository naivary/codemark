package converter

import (
	"reflect"

	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/marker"
)

var _ converter.Converter = (*anyConverter)(nil)

type anyConverter struct {
	name string
}

func Any() converter.Converter {
	return &anyConverter{
		name: "any",
	}
}

func (a *anyConverter) Name() string {
	return converter.NewName(_codemark, a.name)
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
	return converter.ConvertTo(m.Value, to)
}
