package converter

import (
	"fmt"
	"reflect"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	"github.com/naivary/codemark/marker"
)

var _ convv1.Converter = (*listConverter)(nil)

type listConverter struct {
	name string
	mngr *Manager
}

func NewList(mngr *Manager) convv1.Converter {
	return &listConverter{
		name: "list",
		mngr: mngr,
	}
}

func (l *listConverter) Name() string {
	return NewName(_codemark, l.name)
}

func (l *listConverter) SupportedTypes() []reflect.Type {
	types := []any{
		// any
		[]any{},
		// int
		[]int{},
		[]int8{},
		[]int16{},
		// =[]rune
		[]int32{},
		[]int64{},

		// uint
		[]uint{},
		// =[]byte
		[]uint8{},
		[]uint16{},
		[]uint32{},
		[]uint64{},

		// float
		[]float32{},
		[]float64{},

		// complex
		[]complex64{},
		[]complex128{},

		// singles
		[]string{},
		[]bool{},

		// ptr int
		[]*int{},
		[]*int8{},
		[]*int16{},
		// []*rune
		[]*int32{},
		[]*int64{},

		// ptr uint
		[]*uint{},
		// []*byte
		[]*uint8{},
		[]*uint16{},
		[]*uint32{},
		[]*uint64{},

		// ptr float
		[]*float32{},
		[]*float64{},

		// ptr complex
		[]*complex64{},
		[]*complex128{},

		// ptr singles
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

func (l *listConverter) CanConvert(m marker.Marker, to reflect.Type) error {
	if m.Kind != marker.LIST {
		return fmt.Errorf(
			"marker kind of `%s` cannot be converted to a string. valid option is: %s",
			m.Kind,
			marker.LIST,
		)
	}
	return nil
}

func (l *listConverter) Convert(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	list := reflect.New(to).Elem()
	elemType := to.Elem()
	elems := m.Value.Interface().([]any)
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
	conv, err := l.mngr.Get(typ)
	if err != nil {
		return _rvzero, err
	}
	mkind := marker.KindOf(rvalue.Type())
	fakeMarker := marker.Fake(mkind, rvalue)
	return conv.Convert(fakeMarker, typ)
}
