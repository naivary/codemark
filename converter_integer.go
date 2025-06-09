package codemark

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
)

var _ Converter = (*intConverter)(nil)

type intConverter struct{}

func (i *intConverter) SupportedTypes() []reflect.Type {
	types := []any{
		int(0),
		int8(0),
		int16(0),
		// byte
		int32(0),
		int64(0),
		uint(0),
		//rune
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		//pointer
		new(int),
		new(int8),
		new(int16),
		// *byte
		new(int32),
		new(int64),
		new(uint),
		// *rune
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

func (i *intConverter) CanConvert(m parser.Marker, def *Definition) error {
	k := m.Kind()
	out := def.output
	if out.Kind() == reflect.Pointer {
		out = def.output.Elem()
	}
	if k == parser.MarkerKindInt {
		return nil
	}
	if k == parser.MarkerKindString && out.Kind() == _byte || out.Kind() == _rune {
		return nil
	}
	return fmt.Errorf("marker kind of `%s` cannot be converted to a int. valid option is: %s\n", m.Kind(), parser.MarkerKindInt)
}

func (i *intConverter) Convert(m parser.Marker, def *Definition) (reflect.Value, error) {
	// TODO: find a better way instead of the if chains
	typeID := TypeID(def.output)
	switch typeID {
	case TypeIDFromAny(int(0)), TypeIDFromAny(int8(0)), TypeIDFromAny(int16(0)), TypeIDFromAny(int64(0)):
		return i.integer(m, def, false)
	case TypeIDFromAny(new(int)), TypeIDFromAny(new(int8)), TypeIDFromAny(new(int16)), TypeIDFromAny(new(int64)):
		return i.integer(m, def, true)
	case TypeIDFromAny(uint(0)), TypeIDFromAny(uint16(0)), TypeIDFromAny(uint32(0)), TypeIDFromAny(uint64(0)):
		return i.uinteger(m, def, false)
	case TypeIDFromAny(new(uint)), TypeIDFromAny(new(uint16)), TypeIDFromAny(new(uint32)), TypeIDFromAny(new(uint64)):
		return i.uinteger(m, def, true)
	}

	if m.Kind() == parser.MarkerKindInt && typeID == TypeIDFromAny(int32(0)) {
		return i.integer(m, def, false)
	}
	if m.Kind() == parser.MarkerKindInt && typeID == TypeIDFromAny(new(int32)) {
		return i.integer(m, def, true)
	}
	if m.Kind() == parser.MarkerKindInt && typeID == TypeIDFromAny(uint8(0)) {
		return i.uinteger(m, def, false)
	}
	if m.Kind() == parser.MarkerKindInt && typeID == TypeIDFromAny(new(uint8)) {
		return i.uinteger(m, def, true)
	}

	if m.Kind() == parser.MarkerKindString && typeID == TypeIDFromAny(int32(0)) {
		return i.bytee(m, def, false)
	}
	if m.Kind() == parser.MarkerKindString && typeID == TypeIDFromAny(new(int32)) {
		return i.bytee(m, def, true)
	}
	if m.Kind() == parser.MarkerKindString && typeID == TypeIDFromAny(uint8(0)) {
		return i.runee(m, def, false)
	}
	if m.Kind() == parser.MarkerKindString && typeID == TypeIDFromAny(new(uint8)) {
		return i.runee(m, def, true)
	}

	return _rvzero, fmt.Errorf("conversion of `%s` to `%s` is not possible", m.Ident(), def.output)
}

func (i *intConverter) integer(m parser.Marker, def *Definition, isPtr bool) (reflect.Value, error) {
	n := m.Value().Int()
	if i.isOverflowingInt(def.output, n) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`\n", m, def.output)
	}
	return toOutput(m.Value(), def.output, isPtr)
}

func (i *intConverter) uinteger(m parser.Marker, def *Definition, isPtr bool) (reflect.Value, error) {
	n := m.Value().Int()
	if i.isOverflowingUint(def.output, uint64(n)) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`\n", m, def.output)
	}
	return toOutput(m.Value(), def.output, isPtr)
}

func (i *intConverter) runee(m parser.Marker, def *Definition, isPtr bool) (reflect.Value, error) {
	markerValue := m.Value().String()
	if len(markerValue) > 1 {
		return _rvzero, fmt.Errorf("marker value cannot be bigger than 2 chars for rune conversion: %s\n", m.Value())
	}
	rvalue := reflect.ValueOf(rune(markerValue[0]))
	return toOutput(rvalue, def.output, isPtr)
}

func (i *intConverter) bytee(m parser.Marker, def *Definition, isPtr bool) (reflect.Value, error) {
	markerValue := m.Value().String()
	if len(markerValue) > 1 {
		return _rvzero, fmt.Errorf("value of marker is bigger than 2: %s\n", m.Value())
	}
	bvalue := reflect.ValueOf(byte(markerValue[0]))
	return toOutput(bvalue, def.output, isPtr)
}

func (i *intConverter) isOverflowingInt(out reflect.Type, n int64) bool {
	if out.Kind() == reflect.Pointer {
		out = out.Elem()
	}
	return out.OverflowInt(n)
}

func (i *intConverter) isOverflowingUint(out reflect.Type, n uint64) bool {
	if out.Kind() == reflect.Pointer {
		out = out.Elem()
	}
	return out.OverflowUint(n)
}
