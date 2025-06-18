package codemark

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Converter = (*boolConverter)(nil)

type boolConverter struct{}

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

func (b *boolConverter) CanConvert(m sdk.Marker, def *sdk.Definition) error {
	if m.Kind() != sdk.MarkerKindBool {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a boolean. valid option is: %s\n", m.Kind(), sdk.MarkerKindBool)
	}
	return nil
}

func (b *boolConverter) Convert(m sdk.Marker, def *sdk.Definition) (reflect.Value, error) {
	return sdkutil.ConvertTo(m.Value(), def.Output)
}
