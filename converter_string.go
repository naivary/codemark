package codemark

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
)

var _ Converter = (*stringConverter)(nil)

type stringConverter struct{}

func (s *stringConverter) SupportedTypes() []reflect.Type {
	types := []any{
		string(""),
		rune(0),
		byte(0),
		[]byte{},
		// pointer
		new(string),
		new(rune),
		new(byte),
		[]*byte{},
	}
	supported := make([]reflect.Type, 0, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		supported = append(supported, rtype)
	}
	return supported
}

func (s *stringConverter) CanConvert(m parser.Marker, def *Definition) error {
	if m.Kind() != parser.MarkerKindString {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a string. valid option is: %s\n", m.Kind(), parser.MarkerKindString)
	}
	// NOTE: dont need to check if the def.output is supported because the converter
	// will only be choosen if def.output is one of the supported types
	return nil
}

func (s *stringConverter) Convert(m parser.Marker, def *Definition) (any, error) {
	typeID := TypeID(def.output)
	switch typeID {
	case TypeIDFromAny(string("")):
		return s.str(m, def, false)
	case TypeIDFromAny(new(string)):
		return s.str(m, def, true)
	case TypeIDFromAny(rune(0)):
		return s.runee(m, def, false)
	case TypeIDFromAny(new(rune)):
		return s.runee(m, def, true)
	case TypeIDFromAny(byte(0)):
		return s.bytee(m, def, false)
	case TypeIDFromAny(new(byte)):
		return s.bytee(m, def, true)
	case TypeIDFromAny([]byte{}):
		return s.bytes(m, def, false)
	case TypeIDFromAny([]*byte{}):
		return s.bytes(m, def, true)
	}
	return nil, fmt.Errorf("conversion of `%s` to `%s` is not possible", m.Ident(), def.output)
}

func (s *stringConverter) str(m parser.Marker, def *Definition, isPtr bool) (any, error) {
	typ := reflect.TypeOf("")
	v := reflect.New(typ)
	v.Elem().SetString(m.Value().String())
	return toOutput(v, def, isPtr)
}

func (s *stringConverter) runee(m parser.Marker, def *Definition, isPtr bool) (any, error) {
	markerValue := m.Value().String()
	if len(markerValue) > 1 {
		return nil, fmt.Errorf("marker value cannot be bigger than 2 chars for rune conversion: %s\n", m.Value())
	}
	typ := reflect.TypeOf(rune(0))
	r := reflect.New(typ)
	rvalue := reflect.ValueOf(rune(markerValue[0]))
	r.Elem().Set(rvalue)
	return toOutput(r, def, isPtr)
}

func (s *stringConverter) bytee(m parser.Marker, def *Definition, isPtr bool) (any, error) {
	markerValue := m.Value().String()
	if len(markerValue) > 1 {
		return nil, fmt.Errorf("value of marker is bigger than 2: %s\n", m.Value())
	}
	typ := reflect.TypeOf(byte(0))
	b := reflect.New(typ)
	bvalue := reflect.ValueOf(byte(markerValue[0]))
	b.Elem().Set(bvalue)
	return toOutput(b, def, isPtr)
}

func (s *stringConverter) bytes(m parser.Marker, def *Definition, isPtr bool) (any, error) {
	markerValue := m.Value().String()
	elem := reflect.TypeOf(byte(0))
	slice := reflect.TypeOf([]byte{})
	if isPtr {
		slice = reflect.TypeOf([]*byte{})
	}
	bytes := reflect.MakeSlice(slice, 0, len(markerValue))
	for _, b := range []byte(markerValue) {
		rbyte := reflect.New(elem)
		rvalue := reflect.ValueOf(b)
		rbyte.Elem().Set(rvalue)
		if isPtr {
			bytes = reflect.Append(bytes, rbyte)
			continue
		}
		bytes = reflect.Append(bytes, rbyte.Elem())
	}
	// isPtr is always true because you cannot dereference a slice
	return toOutput(bytes, def, true)
}
