package converter

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/typeutil"
)

var _ converter.Converter = (*complexConverter)(nil)

type complexConverter struct {
	name string
}

func Complex() converter.Converter {
	return &complexConverter{
		name: "complex",
	}
}

func (c *complexConverter) Name() string {
	return converter.NewName(_codemark, c.name)
}

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

func (c *complexConverter) CanConvert(m marker.Marker, to reflect.Type) error {
	if m.Kind != marker.COMPLEX {
		return fmt.Errorf(
			"marker kind of `%s` cannot be converted to a string. valid option is: %s",
			m.Kind,
			marker.COMPLEX,
		)
	}
	return nil
}

func (c *complexConverter) Convert(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	return c.complexx(m, to)
}

func (c *complexConverter) complexx(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	n := m.Value.Complex()
	if c.isOverflowing(to, n) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`", m.String(), to)
	}
	return converter.ConvertTo(m.Value, to)
}

func (c *complexConverter) isOverflowing(out reflect.Type, n complex128) bool {
	return typeutil.Deref(out).OverflowComplex(n)
}
