package converter

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Converter = (*intConverter)(nil)

type intConverter struct {
	name string
}

func Integer() sdk.Converter {
	return &intConverter{
		name: "integer",
	}
}

func (i *intConverter) Name() string {
	return sdkutil.NewConvName(_codemark, i.name)
}

func (i *intConverter) SupportedTypes() []reflect.Type {
	types := []any{
		int(0),
		int8(0),
		int16(0),
		// rune
		int32(0),
		int64(0),
		uint(0),
		// byte
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		// pointer
		new(int),
		new(int8),
		new(int16),
		// *rune
		new(int32),
		new(int64),
		new(uint),
		// *byte
		new(uint8),
		new(uint16),
		new(uint32),
		new(uint64),
	}
	supported := make([]reflect.Type, 0, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		supported = append(supported, rtype)
	}
	return supported
}

func (i *intConverter) CanConvert(m parser.Marker, def *sdk.Definition) error {
	mkind := m.Kind()
	out := sdkutil.Deref(def.Output)
	if mkind == marker.INT {
		return nil
	}
	if mkind == marker.STRING && out.Kind() == _byte || out.Kind() == _rune {
		return nil
	}
	return fmt.Errorf("marker kind of `%s` cannot be converted to a int. valid options are: %s;%s\n", m.Kind(), marker.INT, marker.STRING)
}

func (i *intConverter) Convert(m parser.Marker, def *sdk.Definition) (reflect.Value, error) {
	mkind := m.Kind()
	if i.isInteger(def.Output, mkind) {
		return i.integer(m, def)
	}
	if i.isUint(def.Output, mkind) {
		return i.uinteger(m, def)
	}
	if i.isByte(def.Output, mkind) {
		return i.bytee(m, def)
	}
	if i.isRune(def.Output, mkind) {
		return i.runee(m, def)
	}
	return _rvzero, fmt.Errorf("cannot converter %s to %v\n", m.Ident(), def.Output)
}

func (i *intConverter) integer(m parser.Marker, def *sdk.Definition) (reflect.Value, error) {
	n := m.Value().Int()
	if i.isOverflowingInt(def.Output, n) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`\n", m.String(), def.Output)
	}
	return sdkutil.ConvertTo(m.Value(), def.Output)
}

func (i *intConverter) uinteger(m parser.Marker, def *sdk.Definition) (reflect.Value, error) {
	n := m.Value().Int()
	if i.isOverflowingUint(def.Output, uint64(n)) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`\n", m.String(), def.Output)
	}
	return sdkutil.ConvertTo(m.Value(), def.Output)
}

func (i *intConverter) runee(m parser.Marker, def *sdk.Definition) (reflect.Value, error) {
	markerValue := m.Value().String()
	if len(markerValue) > 1 {
		return _rvzero, fmt.Errorf("marker value cannot be bigger than 2 chars for rune conversion: %s\n", m.Value())
	}
	rvalue := reflect.ValueOf(rune(markerValue[0]))
	return sdkutil.ConvertTo(rvalue, def.Output)
}

func (i *intConverter) bytee(m parser.Marker, def *sdk.Definition) (reflect.Value, error) {
	markerValue := m.Value().String()
	if len(markerValue) > 1 {
		return _rvzero, fmt.Errorf("value of marker is bigger than 2: %s\n", m.Value())
	}
	bvalue := reflect.ValueOf(byte(markerValue[0]))
	return sdkutil.ConvertTo(bvalue, def.Output)
}

func (i *intConverter) isOverflowingInt(out reflect.Type, n int64) bool {
	return sdkutil.Deref(out).OverflowInt(n)
}

func (i *intConverter) isOverflowingUint(out reflect.Type, n uint64) bool {
	return sdkutil.Deref(out).OverflowUint(n)
}

func (i *intConverter) isInteger(rtype reflect.Type, mkind marker.Kind) bool {
	return sdkutil.IsInt(rtype) && mkind == marker.INT
}

func (i *intConverter) isUint(rtype reflect.Type, mkind marker.Kind) bool {
	return sdkutil.IsUint(rtype) || mkind == marker.INT
}

func (i *intConverter) isByte(rtype reflect.Type, mkind marker.Kind) bool {
	return sdkutil.IsInt(rtype) && mkind == marker.STRING
}

func (i *intConverter) isRune(rtype reflect.Type, mkind marker.Kind) bool {
	return sdkutil.IsUint(rtype) && mkind == marker.STRING
}
