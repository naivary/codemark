package codemark

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Converter = (*boolConverter)(nil)

type boolConverter struct{}

func (b *boolConverter) Name() string {
	return "codemark.bool"
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

func (b *boolConverter) CanConvert(m parser.Marker, def *sdk.Definition) error {
	if m.Kind() != marker.BOOL {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a boolean. valid option is: %s\n", m.Kind(), marker.BOOL)
	}
	return nil
}

func (b *boolConverter) Convert(m parser.Marker, def *sdk.Definition) (reflect.Value, error) {
	return sdkutil.ConvertTo(m.Value(), def.Output)
}
