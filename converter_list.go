package codemark

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
)

var _ Converter = (*listConverter)(nil)

type listConverter struct{}

func (l *listConverter) SupportedTypes() []reflect.Type {
	types := []any{}
	supported := make([]reflect.Type, 0, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		supported = append(supported, rtype)
	}
	return supported
}

func (l *listConverter) CanConvert(m parser.Marker, def *Definition) error {
	if m.Kind() != parser.MarkerKindList {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a string. valid option is: %s\n", m.Kind(), parser.MarkerKindString)
	}
	return nil
}

func (l *listConverter) Convert(m parser.Marker, def *Definition) (any, error) {
	typeID := TypeID(def.output)
	switch typeID {
	}
	return nil, fmt.Errorf("conversion of `%s` to `%s` is not possible", m.Ident(), def.output)
}
