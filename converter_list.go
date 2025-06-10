package codemark

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
)

var _ Converter = (*listConverter)(nil)

type listConverter struct {
	mngr *ConverterManager
}

func (l *listConverter) SupportedTypes() []reflect.Type {
	types := []any{
		// int
		[]int{},
		[]int8{},
		[]int16{},
		// =[]byte
		[]int32{},
		[]int64{},
		// uint
		[]uint{},
		// =[]rune
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

		// pointer
		// int
		[]*int{},
		[]*int8{},
		[]*int16{},
		// []*byte
		[]*int32{},
		[]*int64{},
		// uint
		[]*uint{},
		// []*rune
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
	// TODO: also accept parser.MarkerKindString for []bytes and []rune
	// conversion
	return nil
}

func (l *listConverter) Convert(m parser.Marker, def *Definition) (reflect.Value, error) {
	typeID := TypeID(def.output)
	switch typeID {
	case TypeIDFromAny([]string{}), TypeIDFromAny([]rune{}):
		return l.list(m, def, false)
	case TypeIDFromAny([]*string{}):
		return l.list(m, def, true)
	}
	return _rvzero, fmt.Errorf("conversion of `%s` to `%s` is not possible", m.Ident(), def.output)
}

func (l *listConverter) list(m parser.Marker, def *Definition, isPtr bool) (reflect.Value, error) {
	list := reflect.New(def.output).Elem()
	elemType := def.output.Elem()
	elems := m.Value().Interface().([]any)
	for _, elem := range elems {
		elemValue, err := l.elem(elem, elemType)
		if err != nil {
			return _rvzero, err
		}
		list = reflect.Append(list, elemValue)
	}
	return list, nil
}

func (l *listConverter) elem(v any, typ reflect.Type) (reflect.Value, error) {
	rvalue := reflect.ValueOf(v)
	conv, err := l.mngr.GetConverter(typ)
	if err != nil {
		return _rvzero, err
	}
	// TODO: Need a reflect.Kind to MarkerKind to dynamically asses which
	// parser.MarkerKindString to use and pass the `CanConvert` assertions of
	// the converter
	markerKind := parser.MarkerKindOf(rvalue.Type())
	fakeMarker := parser.NewMarker("", markerKind, rvalue)
	fakeDef := MakeDef("", TargetField, typ)
	return conv.Convert(fakeMarker, fakeDef)
}
