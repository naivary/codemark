package converter

import (
	"reflect"

	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Converter = (*anyConverter)(nil)

type anyConverter struct {
	name string
}

func Any() sdk.Converter {
	return &anyConverter{
		name: "any",
	}
}

func (a *anyConverter) Name() string {
	return sdkutil.NewConvName(_codemark, a.name)
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
	return sdkutil.ConvertTo(m.Value, to)
}
