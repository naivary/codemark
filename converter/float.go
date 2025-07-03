package converter

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/definition"
	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Converter = (*floatConverter)(nil)

type floatConverter struct {
	name string
}

func Float() sdk.Converter {
	return &floatConverter{
		name: "float",
	}
}

func (f *floatConverter) Name() string {
	return sdkutil.NewConvName(_codemark, f.name)
}

func (f *floatConverter) SupportedTypes() []reflect.Type {
	types := []any{
		float32(0.0),
		float64(0.0),
		// pointer
		new(float32),
		new(float64),
	}
	supported := make([]reflect.Type, 0, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		supported = append(supported, rtype)
	}
	return supported
}

func (f *floatConverter) CanConvert(m marker.Marker, def *definition.Definition) error {
	if m.Kind != marker.FLOAT {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a float. valid option is: %s\n", m.Kind, marker.FLOAT)
	}
	return nil
}

func (f *floatConverter) Convert(m marker.Marker, def *definition.Definition) (reflect.Value, error) {
	return f.float(m, def)
}

func (f *floatConverter) float(m marker.Marker, def *definition.Definition) (reflect.Value, error) {
	n := m.Value.Float()
	if f.isOverflowing(def.Output, n) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`\n", m.String(), def.Output)
	}
	return sdkutil.ConvertTo(m.Value, def.Output)
}

func (f *floatConverter) isOverflowing(out reflect.Type, n float64) bool {
	return sdkutil.Deref(out).OverflowFloat(n)
}
