package converter

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/definition"
	"github.com/naivary/codemark/maker"
	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Converter = (*listConverter)(nil)

type listConverter struct {
	mngr sdk.ConverterManager
	name string
}

func List(mngr sdk.ConverterManager) sdk.Converter {
	return &listConverter{
		mngr: mngr,
		name: "list",
	}
}

func (l *listConverter) Name() string {
	return sdkutil.NewConvName(_codemark, l.name)
}

func (l *listConverter) SupportedTypes() []reflect.Type {
	types := []any{
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

		//float
		[]float32{},
		[]float64{},

		//complex
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

func (l *listConverter) CanConvert(m parser.Marker, def *definition.Definition) error {
	if m.Kind() != marker.LIST {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a string. valid option is: %s\n", m.Kind(), marker.LIST)
	}
	return nil
}

func (l *listConverter) Convert(m parser.Marker, def *definition.Definition) (reflect.Value, error) {
	return l.list(m, def)
}

func (l *listConverter) list(m parser.Marker, def *definition.Definition) (reflect.Value, error) {
	list := reflect.New(def.Output).Elem()
	elemType := def.Output.Elem()
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
	mkind := sdkutil.MarkerKindOf(rvalue.Type())
	fakeMarker := maker.MakeFakeMarker(mkind, rvalue)
	fakeDef, err := maker.MakeFakeDef(typ)
	if err != nil {
		return _rvzero, err
	}
	return conv.Convert(fakeMarker, fakeDef)
}
