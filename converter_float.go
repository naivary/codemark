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

func (f *floatConverter) Convert(m parser.Marker, def *Definition) (any, error) {
	typeID := TypeID(def.output)
	switch typeID {
	case TypeIDFromAny(float32(0.0)), TypeIDFromAny(float64(0.0)):
		return f.float(m, def, false)
	case TypeIDFromAny(new(float32)), TypeIDFromAny(new(float64)):
		return f.float(m, def, true)
	}
	return nil, fmt.Errorf("conversion of `%s` to `%s` is not possible", m.Ident(), def.output)
}

func (f *floatConverter) float(m parser.Marker, def *Definition, isPtr bool) (any, error) {
	n := m.Value().Float()
	if f.isOverflowing(def.output, n) {
		return nil, fmt.Errorf("overflow converting `%s` to `%v`\n", m, def.output)
	}
	out, err := toOutput(m.Value(), def.output, isPtr)
	if err != nil {
		return nil, err
	}
	return out.Interface(), nil
}

func (f *floatConverter) isOverflowing(out reflect.Type, n float64) bool {
	if out.Kind() == reflect.Pointer {
		out = out.Elem()
	}
	return out.OverflowFloat(n)
}
