package converter

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/api"
	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Converter = (*complexConverter)(nil)

type complexConverter struct {
	name string
}

func Complex() sdk.Converter {
	return &complexConverter{
		name: "complex",
	}
}

func (c *complexConverter) Name() string {
	return sdkutil.NewConvName(_codemark, c.name)
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

func (c *complexConverter) CanConvert(m marker.Marker, def *api.Definition) error {
	if m.Kind != marker.COMPLEX {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a string. valid option is: %s\n", m.Kind, marker.COMPLEX)
	}
	return nil
}

func (c *complexConverter) Convert(m marker.Marker, def *api.Definition) (reflect.Value, error) {
	return c.complexx(m, def)
}

func (c *complexConverter) complexx(m marker.Marker, def *api.Definition) (reflect.Value, error) {
	n := m.Value.Complex()
	if c.isOverflowing(def.Output, n) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`\n", m.String(), def.Output)
	}
	return sdkutil.ConvertTo(m.Value, def.Output)
}

func (c *complexConverter) isOverflowing(out reflect.Type, n complex128) bool {
	return sdkutil.Deref(out).OverflowComplex(n)
}
