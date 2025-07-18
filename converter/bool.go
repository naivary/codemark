package converter

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Converter = (*boolConverter)(nil)

type boolConverter struct{}

func Bool() sdk.Converter {
	return &boolConverter{}
}

func (b *boolConverter) Name() string {
	return sdkutil.NewConvName(_codemark, "bool")
}

func (b *boolConverter) SupportedTypes() []reflect.Type {
	types := []any{
		bool(false),
		//pointer
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
		return fmt.Errorf("marker kind of `%s` cannot be converted to a boolean. valid option is: %s", m.Kind, marker.BOOL)
	}
	return nil
}

func (b *boolConverter) Convert(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	return sdkutil.ConvertTo(m.Value, to)
}
