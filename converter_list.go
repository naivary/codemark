package codemark

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
)

var _ Converter = (*listConverter)(nil)

type listConverter struct{}

func (l *listConverter) SupportedTypes() []reflect.Type {
	types := []any{
		// int
		[]int{},
		[]int8{},
		[]int16{},
		[]int32{},
		[]int64{},
		// uint
		[]uint{},
		[]uint8{},
		[]uint16{},
		[]uint32{},
		[]uint64{},
		//float
		[]float32{},
		[]float64{},
		//complex
		[]complex64{},
		[]complex128{},
		// singles
		[]string{},
		[]bool{},
		[]byte{},
		[]rune{},

		// pointer
		// int
		[]*int{},
		[]*int8{},
		[]*int16{},
		[]*int32{},
		[]*int64{},
		// uint
		[]*uint{},
		[]*uint8{},
		[]*uint16{},
		[]*uint32{},
		[]*uint64{},
		// float
		[]*float32{},
		[]*float64{},
		// complex
		[]*complex64{},
		[]*complex128{},
		// singles
		[]*string{},
		[]*bool{},
		[]*byte{},
		[]*rune{},
	}
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

func (l *listConverter) Convert(m parser.Marker, def *Definition) (reflect.Value, error) {
	typeID := TypeID(def.output)
	switch typeID {
	case TypeIDFromAny([]string{}):
		return l.str(m, def, false)
	case TypeIDFromAny([]*string{}):
		return l.str(m, def, true)
	}
	return _rvzero, fmt.Errorf("conversion of `%s` to `%s` is not possible", m.Ident(), def.output)
}

func (l *listConverter) str(m parser.Marker, def *Definition, isPtr bool) (reflect.Value, error) {
	list := reflect.New(def.output).Elem()
	elemType := def.output.Elem()
	elems := m.Value().Interface().([]any)
	for _, elem := range elems {
		rvalue := reflect.ValueOf(elem)
		out, err := toOutput(rvalue, elemType, isPtr)
		if err != nil {
			return _rvzero, err
		}
		list = reflect.Append(list, out)
	}
	return list, nil
}
