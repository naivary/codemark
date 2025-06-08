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
	out, err := toOutput(m.Value(), def.output, isPtr)
	if err != nil {
		return nil, err
	}
	return out.Interface(), nil
}

func (s *stringConverter) runee(m parser.Marker, def *Definition, isPtr bool) (any, error) {
	markerValue := m.Value().String()
	if len(markerValue) > 1 {
		return nil, fmt.Errorf("marker value cannot be bigger than 2 chars for rune conversion: %s\n", m.Value())
	}
	rvalue := reflect.ValueOf(rune(markerValue[0]))
	out, err := toOutput(rvalue, def.output, isPtr)
	if err != nil {
		return nil, err
	}
	return out.Interface(), nil

}

func (s *stringConverter) bytee(m parser.Marker, def *Definition, isPtr bool) (any, error) {
	markerValue := m.Value().String()
	if len(markerValue) > 1 {
		return nil, fmt.Errorf("value of marker is bigger than 2: %s\n", m.Value())
	}
	bvalue := reflect.ValueOf(byte(markerValue[0]))
	out, err := toOutput(bvalue, def.output, isPtr)
	if err != nil {
		return nil, err
	}
	return out.Interface(), nil
}

func (s *stringConverter) bytes(m parser.Marker, def *Definition, isPtr bool) (any, error) {
	v := m.Value()
	bytes := reflect.MakeSlice(def.output, 0, len(v.String()))
	elemType := def.output.Elem()
	for _, b := range []byte(v.String()) {
		rvalue := reflect.ValueOf(b)
		elem, err := toOutput(rvalue, elemType, isPtr)
		if err != nil {
			return nil, err
		}
		bytes = reflect.Append(bytes, elem)
	}
	return bytes.Interface(), nil
}
