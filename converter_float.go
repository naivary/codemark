package codemark

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
)

var _ Converter = (*floatConverter)(nil)

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

func (f *floatConverter) CanConvert(m parser.Marker, def *Definition) error {
	if m.Kind() != parser.MarkerKindFloat {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a float. valid option is: %s\n", m.Kind(), parser.MarkerKindFloat)
	}
	return nil
}

func (f *floatConverter) Convert(m parser.Marker, def *Definition) (reflect.Value, error) {
	return f.float(m, def)
}

func (f *floatConverter) float(m parser.Marker, def *Definition) (reflect.Value, error) {
	n := m.Value().Float()
	if f.isOverflowing(def.output, n) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`\n", m, def.output)
	}
	return toOutput(m.Value(), def.output)
}

func (f *floatConverter) isOverflowing(out reflect.Type, n float64) bool {
	if out.Kind() == reflect.Pointer {
		out = out.Elem()
	}
	return out.OverflowFloat(n)
}
