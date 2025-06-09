package codemark

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
)

var _ Converter = (*boolConverter)(nil)

type boolConverter struct{}

func (b *boolConverter) SupportedTypes() []reflect.Type {
	types := []any{
		bool(false),
		//pointer
		new(bool),
	}
	supported := make([]reflect.Type, 0, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		supported = append(supported, rtype)
	}
	return supported
}

func (b *boolConverter) CanConvert(m parser.Marker, def *Definition) error {
	if m.Kind() != parser.MarkerKindBool {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a boolean. valid option is: %s\n", m.Kind(), parser.MarkerKindBool)
	}
	return nil
}

func (b *boolConverter) Convert(m parser.Marker, def *Definition) (any, error) {
	typeID := TypeID(def.output)
	switch typeID {
	case TypeIDFromAny(bool(false)):
		return b.boolean(m, def, false)
	case TypeIDFromAny(new(bool)):
		return b.boolean(m, def, true)
	}
	return nil, fmt.Errorf("conversion of `%s` to `%s` is not possible", m.Ident(), def.output)
}

func (b *boolConverter) boolean(m parser.Marker, def *Definition, isPtr bool) (any, error) {
	out, err := toOutput(m.Value(), def.output, isPtr)
	if err != nil {
		return nil, err
	}
	return out.Interface(), nil
}
