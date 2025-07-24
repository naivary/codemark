package converter

import (
	"fmt"
	"reflect"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/typeutil"
)

var _ convv1.Converter = (*floatConverter)(nil)

type floatConverter struct {
	name string
}

func NewFloat() convv1.Converter {
	return &floatConverter{
		name: "float",
	}
}

func (f *floatConverter) Name() string {
	return NewName(_codemark, f.name)
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

func (f *floatConverter) CanConvert(m marker.Marker, to reflect.Type) error {
	if m.Kind != marker.FLOAT {
		return fmt.Errorf(
			"marker kind of `%s` cannot be converted to a float. valid option is: %s",
			m.Kind,
			marker.FLOAT,
		)
	}
	return nil
}

func (f *floatConverter) Convert(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	n := m.Value.Float()
	if f.isOverflowing(to, n) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`", m.String(), to)
	}
	return ConvertTo(m.Value, to)
}

func (f *floatConverter) isOverflowing(out reflect.Type, n float64) bool {
	return typeutil.Deref(out).OverflowFloat(n)
}
