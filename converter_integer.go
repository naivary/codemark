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
		int32(0),
		int64(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		//pointer
		new(int),
		new(int8),
		new(int16),
		new(int32),
		new(int64),
		new(uint),
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
	if m.Kind() != parser.MarkerKindInt {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a int. valid option is: %s\n", m.Kind(), parser.MarkerKindInt)
	}
	return nil
}

func (i *intConverter) Convert(m parser.Marker, def *Definition) (any, error) {
	typeID := TypeID(def.output)
	switch typeID {
	case TypeIDFromAny(int(0)), TypeIDFromAny(int8(0)), TypeIDFromAny(int16(0)), TypeIDFromAny(int32(0)), TypeIDFromAny(int64(0)):
		return i.integer(m, def, false)
	case TypeIDFromAny(new(int)), TypeIDFromAny(new(int8)), TypeIDFromAny(new(int16)), TypeIDFromAny(new(int32)), TypeIDFromAny(new(int64)):
		return i.integer(m, def, true)
	case TypeIDFromAny(uint(0)), TypeIDFromAny(uint8(0)), TypeIDFromAny(uint16(0)), TypeIDFromAny(uint32(0)), TypeIDFromAny(uint64(0)):
		return i.uinteger(m, def, false)
	case TypeIDFromAny(new(uint)), TypeIDFromAny(new(uint8)), TypeIDFromAny(new(uint16)), TypeIDFromAny(new(uint32)), TypeIDFromAny(new(uint64)):
		return i.uinteger(m, def, true)
	}
	return nil, fmt.Errorf("conversion of `%s` to `%s` is not possible", m.Ident(), def.output)
}

func (i *intConverter) integer(m parser.Marker, def *Definition, isPtr bool) (any, error) {
	n := m.Value().Int()
	if i.isOverflowing(def.output, n) {
		return nil, fmt.Errorf("overflow converting `%s` to `%v`\n", m, def.output)
	}

	outputType := def.output
	if isPtr {
		outputType = def.output.Elem()
	}
	out := reflect.New(outputType)
	out.Elem().SetInt(n)
	if isPtr {
		return out.Convert(def.output).Interface(), nil
	}
	return out.Elem().Interface(), nil
}

func (i *intConverter) uinteger(m parser.Marker, def *Definition, isPtr bool) (any, error) {
	return nil, nil
}

func (i *intConverter) isOverflowing(out reflect.Type, n int64) bool {
	if out.Kind() == reflect.Pointer {
		out = out.Elem()
	}
	return out.OverflowInt(n)
}
