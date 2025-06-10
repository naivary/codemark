package codemark

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
)

var _ Converter = (*complexConverter)(nil)

type complexConverter struct{}

func (c *complexConverter) SupportedTypes() []reflect.Type {
	types := []any{
		complex64(0 + 0i),
		complex128(0 + 0i),
		// pointer
		new(complex64),
		new(complex128),
	}
	supported := make([]reflect.Type, 0, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		supported = append(supported, rtype)
	}
	return supported
}

func (c *complexConverter) CanConvert(m parser.Marker, def *Definition) error {
	if m.Kind() != parser.MarkerKindComplex {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a string. valid option is: %s\n", m.Kind(), parser.MarkerKindComplex)
	}
	return nil
}

func (c *complexConverter) Convert(m parser.Marker, def *Definition) (reflect.Value, error) {
	return c.complexx(m, def)
}

func (c *complexConverter) complexx(m parser.Marker, def *Definition) (reflect.Value, error) {
	n := m.Value().Complex()
	if c.isOverflowing(def.output, n) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`\n", m, def.output)
	}
	return toOutput(m.Value(), def.output)
}

func (c *complexConverter) isOverflowing(out reflect.Type, n complex128) bool {
	if out.Kind() == reflect.Pointer {
		out = out.Elem()
	}
	return out.OverflowComplex(n)
}
