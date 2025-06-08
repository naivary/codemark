package codemark

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
)

var _ Converter = (*listConverter)(nil)

type listConverter struct{}

func (i *listConverter) SupportedTypes() []reflect.Type {
	types := []any{
		[]string{},
		[]int{},
		[]uint{},
		[]complex64{},
		[]complex128{},
		[]float32{},
		[]float64{},
	}
	supported := make([]reflect.Type, 0, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		supported = append(supported, rtype)
	}
	return supported
}

func (i *listConverter) CanConvert(m parser.Marker, def *Definition) error {
	if m.Kind() != parser.MarkerKindList {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a int. valid option is: %s\n", m.Kind(), parser.MarkerKindInt)
	}
	// NOTE: dont need to check if the def.output is supported because the converter
	// will only be choosen if def.output is one of the supported types
	return nil
}

func (i *listConverter) Convert(m parser.Marker, def *Definition) (any, error) {
	typeID := TypeID(def.output)
	switch typeID {
	}
	return nil, errors.New("conversion was not possible")
}
