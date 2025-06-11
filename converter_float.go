package codemark

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Converter = (*floatConverter)(nil)

type floatConverter struct{}

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

func (f *floatConverter) CanConvert(m parser.Marker, def *sdk.Definition) error {
	if m.Kind() != parser.MarkerKindFloat {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a float. valid option is: %s\n", m.Kind(), parser.MarkerKindFloat)
	}
	return nil
}

func (f *floatConverter) Convert(m parser.Marker, def *sdk.Definition) (reflect.Value, error) {
	return f.float(m, def)
}

func (f *floatConverter) float(m parser.Marker, def *sdk.Definition) (reflect.Value, error) {
	n := m.Value().Float()
	if f.isOverflowing(def.Output, n) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`\n", m, def.Output)
	}
	return sdkutil.ToType(m.Value(), def.Output)
}

func (f *floatConverter) isOverflowing(out reflect.Type, n float64) bool {
	return sdkutil.Deref(out).OverflowFloat(n)
}
