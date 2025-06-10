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

func (b *boolConverter) Convert(m parser.Marker, def *Definition) (reflect.Value, error) {
	return b.boolean(m, def)
}

func (b *boolConverter) boolean(m parser.Marker, def *Definition) (reflect.Value, error) {
	return toOutput(m.Value(), def.output)
}
