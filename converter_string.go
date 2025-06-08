package codemark

import (
	"errors"
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
	typeID, err := TypeID(def.output)
	if err != nil {
		return nil, err
	}
	markerValue := m.Value()

	switch typeID {
	case "string":
		return toOutput(markerValue, def)
	case _rune.String():
		return s.char(m, def)
	case "ptr.string":
		return s.ptrStr(m, def)
	}
	return nil, errors.New("conversionw as not possible")
}

func (s *stringConverter) ptrStr(m parser.Marker, def *Definition) (any, error) {
	ptr := reflect.ValueOf(new(string))
	val := ptr.Elem()
	val.SetString(m.Value().String())
	return toOutput(ptr, def)
}

// char is named char because rune is taken
func (s *stringConverter) char(m parser.Marker, def *Definition) (any, error) {
	markerValue := m.Value().String()
	if len(markerValue) != 1 {
		return nil, fmt.Errorf("lenght of marker value is unequal to one, making the conversion to a rune impossible: %v\n", markerValue)
	}
	r := reflect.ValueOf(rune(0))
	markerValueRune := reflect.ValueOf(rune(markerValue[0]))
	r.Set(markerValueRune)
	return toOutput(r, def)
}
