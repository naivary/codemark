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
		// pointer
		new(string),
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

func (s *stringConverter) Convert(m parser.Marker, def *Definition) (reflect.Value, error) {
	typeID := TypeID(def.output)
	switch typeID {
	case TypeIDFromAny(string("")):
		return s.str(m, def, false)
	case TypeIDFromAny(new(string)):
		return s.str(m, def, true)
	}
	return _rvzero, fmt.Errorf("conversion of `%s` to `%s` is not possible", m.Ident(), def.output)
}

func (s *stringConverter) str(m parser.Marker, def *Definition, isPtr bool) (reflect.Value, error) {
	return toOutput(m.Value(), def.output, isPtr)
}


