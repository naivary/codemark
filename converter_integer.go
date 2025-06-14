package codemark

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Converter = (*intConverter)(nil)

type intConverter struct{}

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
	if mkind == parser.MarkerKindInt {
		return nil
	}
	if mkind == parser.MarkerKindString && out.Kind() == _byte || out.Kind() == _rune {
		return nil
	}
	return fmt.Errorf("marker kind of `%s` cannot be converted to a int. valid options are: %s;%s\n", m.Kind(), parser.MarkerKindInt, parser.MarkerKindString)
}

func (i *intConverter) Convert(m parser.Marker, def *sdk.Definition) (reflect.Value, error) {
	typeID := sdkutil.TypeIDOf(def.Output)
	markerKind := m.Kind()
	if i.isInteger(typeID, markerKind) {
		return i.integer(m, def)
	}
	if i.isUint(typeID, markerKind) {
		return i.uinteger(m, def)
	}
	if i.isByte(typeID, markerKind) {
		return i.bytee(m, def)
	}
	return i.runee(m, def)
}

func (i *intConverter) integer(m parser.Marker, def *sdk.Definition) (reflect.Value, error) {
	n := m.Value().Int()
	if i.isOverflowingInt(def.Output, n) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`\n", m, def.Output)
	}
	return sdkutil.ConvertTo(m.Value(), def.Output)
}

func (i *intConverter) uinteger(m parser.Marker, def *sdk.Definition) (reflect.Value, error) {
	n := m.Value().Int()
	if i.isOverflowingUint(def.Output, uint64(n)) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`\n", m, def.Output)
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

func (i *intConverter) isInteger(typeID string, mkind parser.MarkerKind) bool {
	return (strings.HasPrefix(typeID, "int") || strings.HasPrefix(typeID, "ptr.int")) && mkind == parser.MarkerKindInt
}

func (i *intConverter) isUint(typeID string, mkind parser.MarkerKind) bool {
	return (strings.HasPrefix(typeID, "uint") || strings.HasPrefix(typeID, "ptr.uint")) && mkind == parser.MarkerKindInt
}

func (i *intConverter) isByte(typeID string, mkind parser.MarkerKind) bool {
	return (strings.HasPrefix(typeID, "ptr.int32") || strings.HasPrefix(typeID, "int32")) && mkind == parser.MarkerKindString
}
