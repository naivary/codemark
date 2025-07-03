package converter

import (
	"reflect"

	"github.com/naivary/codemark/definition"
	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Converter = (*stringConverter)(nil)

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

func (a *anyConverter) CanConvert(m marker.Marker, def *definition.Definition) error {
	return nil
}

func (a *anyConverter) Convert(m marker.Marker, def *definition.Definition) (reflect.Value, error) {
	return a.anything(m, def)
}

func (a *anyConverter) anything(m marker.Marker, def *definition.Definition) (reflect.Value, error) {
	return sdkutil.ConvertTo(m.Value, def.Output)
}
