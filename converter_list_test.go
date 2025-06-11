package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
)

// string
type strList []string

// int
type intList []int
type i8List []int8
type i16List []int16
type byteList []byte // int32
type i64List []int64

// uint
type uintList []uint
type runeList []rune // uint8
type ui16List []uint16
type ui32List []uint32
type ui64List []uint64

// float
type f32List []float32
type f64List []float64

// complex
type c64List []complex64
type c128List []complex128

// bool
type boolList []bool

// ptr string
type ptrstrList []*string

// ptr bool
type ptrboolList []*bool

// ptr int
type ptrintList []*int
type ptri8List []*int8
type ptri16List []*int16
type ptrbyteList []*byte // int32
type ptri64List []*int64

// ptr uint
type ptruintList []*uint
type ptrruneList []*rune // uint8
type ptrui16List []*uint16
type ptrui32List []*uint32
type ptrui64List []*uint64

// float
type ptrf32List []*float32
type ptrf64List []*float64

// complex
type ptrc64List []*complex64
type ptrc128List []*complex128

func listDefs(t *testing.T) sdk.Registry {
	reg := NewInMemoryRegistry()
	defs := []*sdk.Definition{
		// string
		MustMakeDef("path:to:stringlist", sdk.TargetField, reflect.TypeOf(strList([]string{}))),
		// int
		MustMakeDef("path:to:intlist", sdk.TargetField, reflect.TypeOf(intList([]int{}))),
		MustMakeDef("path:to:i8list", sdk.TargetField, reflect.TypeOf(i8List([]int8{}))),
		MustMakeDef("path:to:i16list", sdk.TargetField, reflect.TypeOf(i16List([]int16{}))),
		MustMakeDef("path:to:bytelist", sdk.TargetField, reflect.TypeOf(byteList([]byte{}))), // =in32 list
		MustMakeDef("path:to:i64list", sdk.TargetField, reflect.TypeOf(i64List([]int64{}))),
		// uint
		MustMakeDef("path:to:uintlist", sdk.TargetField, reflect.TypeOf(uintList([]uint{}))),
		MustMakeDef("path:to:runelist", sdk.TargetField, reflect.TypeOf(runeList([]rune{}))), // =uint8 list
		MustMakeDef("path:to:ui16list", sdk.TargetField, reflect.TypeOf(ui16List([]uint16{}))),
		MustMakeDef("path:to:ui32list", sdk.TargetField, reflect.TypeOf(ui32List([]uint32{}))),
		MustMakeDef("path:to:ui64list", sdk.TargetField, reflect.TypeOf(ui64List([]uint64{}))),
		// float
		MustMakeDef("path:to:f32list", sdk.TargetField, reflect.TypeOf(f32List([]float32{}))),
		MustMakeDef("path:to:f64list", sdk.TargetField, reflect.TypeOf(f64List([]float64{}))),
		// complex
		MustMakeDef("path:to:c64list", sdk.TargetField, reflect.TypeOf(c64List([]complex64{}))),
		MustMakeDef("path:to:c128list", sdk.TargetField, reflect.TypeOf(c128List([]complex128{}))),
		// bool
		MustMakeDef("path:to:boollist", sdk.TargetField, reflect.TypeOf(boolList([]bool{}))),

		// ptr string
		MustMakeDef("path:to:ptrstringlist", sdk.TargetField, reflect.TypeOf(ptrstrList([]*string{}))),
		// ptr int
		MustMakeDef("path:to:ptrintlist", sdk.TargetField, reflect.TypeOf(ptrintList([]*int{}))),
		MustMakeDef("path:to:ptri8list", sdk.TargetField, reflect.TypeOf(ptri8List([]*int8{}))),
		MustMakeDef("path:to:ptri16list", sdk.TargetField, reflect.TypeOf(ptri16List([]*int16{}))),
		MustMakeDef("path:to:ptrbytelist", sdk.TargetField, reflect.TypeOf(ptrbyteList([]*byte{}))),
		MustMakeDef("path:to:ptri64list", sdk.TargetField, reflect.TypeOf(ptri64List([]*int64{}))),
		// ptr uint
		MustMakeDef("path:to:ptruintlist", sdk.TargetField, reflect.TypeOf(ptruintList([]*uint{}))),
		MustMakeDef("path:to:ptrrunelist", sdk.TargetField, reflect.TypeOf(ptrruneList([]*rune{}))),
		MustMakeDef("path:to:ptrui16list", sdk.TargetField, reflect.TypeOf(ptrui16List([]*uint16{}))),
		MustMakeDef("path:to:ptrui32list", sdk.TargetField, reflect.TypeOf(ptrui32List([]*uint32{}))),
		MustMakeDef("path:to:ptrui64list", sdk.TargetField, reflect.TypeOf(ptrui64List([]*uint64{}))),
		// ptr float
		MustMakeDef("path:to:ptrf32list", sdk.TargetField, reflect.TypeOf(ptrf32List([]*float32{}))),
		MustMakeDef("path:to:ptrf64list", sdk.TargetField, reflect.TypeOf(ptrf64List([]*float64{}))),
		// ptr complex
		MustMakeDef("path:to:ptrc64list", sdk.TargetField, reflect.TypeOf(ptrc64List([]*complex64{}))),
		MustMakeDef("path:to:ptrc128list", sdk.TargetField, reflect.TypeOf(ptrc128List([]*complex128{}))),
		// ptr bool
		MustMakeDef("path:to:ptrboollist", sdk.TargetField, reflect.TypeOf(ptrboolList([]*bool{}))),
	}
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
			t.Errorf("err occured: %s\n", err)
		}

	}
	return reg
}

func TestListConverter(t *testing.T) {
	tests := []struct {
		name         string
		mrk          parser.Marker
		t            sdk.Target
		out          reflect.Type
		value        []any
		isValid      bool
		isValidValue func(got reflect.Value, expected []any) bool
	}{
		{
			name:    "string list marker to string list",
			mrk:     parser.NewMarker("path:to:stringlist", parser.MarkerKindList, reflect.ValueOf([]any{"path", "to", "marker"})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(strList([]string{})),
			value:   []any{"path", "to", "marker"},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(strList)
				for i, el := range list {
					expectedElem := expected[i].(string)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "string list marker to rune list",
			mrk:     parser.NewMarker("path:to:runelist", parser.MarkerKindList, reflect.ValueOf([]any{"p", "t", "m"})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(runeList([]rune{})),
			value:   []any{rune('p'), rune('t'), rune('m')},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(runeList)
				for i, el := range list {
					expectedElem := expected[i].(rune)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to uint list",
			mrk:     parser.NewMarker("path:to:uintlist", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(uintList([]uint{})),
			value:   []any{uint(3), uint(2), uint(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(uintList)
				for i, el := range list {
					expectedElem := expected[i].(uint)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to rune list",
			mrk:     parser.NewMarker("path:to:runelist", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(runeList([]rune{})),
			value:   []any{rune(3), rune(2), rune(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(runeList)
				for i, el := range list {
					expectedElem := expected[i].(rune)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to uint16 list",
			mrk:     parser.NewMarker("path:to:ui16list", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ui16List([]uint16{})),
			value:   []any{uint16(3), uint16(2), uint16(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ui16List)
				for i, el := range list {
					expectedElem := expected[i].(uint16)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to uint32 list",
			mrk:     parser.NewMarker("path:to:ui32list", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ui32List([]uint32{})),
			value:   []any{uint32(3), uint32(2), uint32(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ui32List)
				for i, el := range list {
					expectedElem := expected[i].(uint32)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to uint64 list",
			mrk:     parser.NewMarker("path:to:ui64list", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ui64List([]uint64{})),
			value:   []any{uint64(3), uint64(2), uint64(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ui64List)
				for i, el := range list {
					expectedElem := expected[i].(uint64)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		// int
		{
			name:    "int list marker to int list",
			mrk:     parser.NewMarker("path:to:intlist", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(intList([]int{})),
			value:   []any{3, 2, 2},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(intList)
				for i, el := range list {
					expectedElem := expected[i].(int)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},

		{
			name:    "int list marker to int8 list",
			mrk:     parser.NewMarker("path:to:i8list", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(i8List([]int8{})),
			value:   []any{int8(3), int8(2), int8(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(i8List)
				for i, el := range list {
					expectedElem := expected[i].(int8)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to int16 list",
			mrk:     parser.NewMarker("path:to:i16list", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(i16List([]int16{})),
			value:   []any{int16(3), int16(2), int16(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(i16List)
				for i, el := range list {
					expectedElem := expected[i].(int16)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "string list marker to byte list",
			mrk:     parser.NewMarker("path:to:bytelist", parser.MarkerKindList, reflect.ValueOf([]any{"p", "t", "m"})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(byteList([]byte{})),
			value:   []any{byte('p'), byte('t'), byte('m')},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(byteList)
				for i, el := range list {
					expectedElem := expected[i].(byte)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to byte list",
			mrk:     parser.NewMarker("path:to:bytelist", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(byteList([]byte{})),
			value:   []any{byte(3), byte(2), byte(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(byteList)
				for i, el := range list {
					expectedElem := expected[i].(byte)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to int64 list",
			mrk:     parser.NewMarker("path:to:i64list", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(i64List([]int64{})),
			value:   []any{int64(3), int64(2), int64(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(i64List)
				for i, el := range list {
					expectedElem := expected[i].(int64)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		// float
		{
			name:    "float list marker to float32 list",
			mrk:     parser.NewMarker("path:to:f32list", parser.MarkerKindList, reflect.ValueOf([]any{3.0, 2.1, 2.2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(f32List([]float32{})),
			value:   []any{float32(3.0), float32(2.1), float32(2.2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(f32List)
				for i, el := range list {
					expectedElem := expected[i].(float32)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "float list marker to float64 list",
			mrk:     parser.NewMarker("path:to:f64list", parser.MarkerKindList, reflect.ValueOf([]any{3.0, 2.1, 2.2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(f64List([]float64{})),
			value:   []any{float64(3.0), float64(2.1), float64(2.2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(f64List)
				for i, el := range list {
					expectedElem := expected[i].(float64)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		// complex
		{
			name:    "complex list marker to complex64 list",
			mrk:     parser.NewMarker("path:to:c64list", parser.MarkerKindList, reflect.ValueOf([]any{0 + 1i, 1 + 2i, 2 + 3i})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(c64List([]complex64{})),
			value:   []any{complex64(0 + 1i), complex64(1 + 2i), complex64(2 + 3i)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(c64List)
				for i, el := range list {
					expectedElem := expected[i].(complex64)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "complex list marker to complex128 list",
			mrk:     parser.NewMarker("path:to:c128list", parser.MarkerKindList, reflect.ValueOf([]any{0 + 1i, 1 + 2i, 2 + 3i})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(c128List([]complex128{})),
			value:   []any{0 + 1i, 1 + 2i, 2 + 3i},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(c128List)
				for i, el := range list {
					expectedElem := expected[i].(complex128)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		// bool
		{
			name:    "bool list marker to bool list",
			mrk:     parser.NewMarker("path:to:boollist", parser.MarkerKindList, reflect.ValueOf([]any{false, true, false})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(boolList([]bool{})),
			value:   []any{false, true, false},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(boolList)
				for i, el := range list {
					expectedElem := expected[i].(bool)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		// pointer
		{
			name:    "string list marker to ptr string list",
			mrk:     parser.NewMarker("path:to:ptrstringlist", parser.MarkerKindList, reflect.ValueOf([]any{"path", "to", "marker"})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrstrList([]*string{})),
			value:   []any{"path", "to", "marker"},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrstrList)
				for i, el := range list {
					expectedElem := expected[i].(string)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to ptr uint list",
			mrk:     parser.NewMarker("path:to:ptruintlist", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptruintList([]*uint{})),
			value:   []any{uint(3), uint(2), uint(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptruintList)
				for i, el := range list {
					expectedElem := expected[i].(uint)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "string list marker to ptr rune list",
			mrk:     parser.NewMarker("path:to:ptrrunelist", parser.MarkerKindList, reflect.ValueOf([]any{"p", "t", "m"})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrruneList([]*rune{})),
			value:   []any{rune('p'), rune('t'), rune('m')},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrruneList)
				for i, el := range list {
					expectedElem := expected[i].(rune)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to ptr rune list",
			mrk:     parser.NewMarker("path:to:ptrrunelist", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrruneList([]*rune{})),
			value:   []any{rune(3), rune(2), rune(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrruneList)
				for i, el := range list {
					expectedElem := expected[i].(rune)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to ptr uint16 list",
			mrk:     parser.NewMarker("path:to:ptrui16list", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrui16List([]*uint16{})),
			value:   []any{uint16(3), uint16(2), uint16(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrui16List)
				for i, el := range list {
					expectedElem := expected[i].(uint16)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to ptr uint32 list",
			mrk:     parser.NewMarker("path:to:ptrui32list", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrui32List([]*uint32{})),
			value:   []any{uint32(3), uint32(2), uint32(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrui32List)
				for i, el := range list {
					expectedElem := expected[i].(uint32)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to ptr uint64 list",
			mrk:     parser.NewMarker("path:to:ptrui64list", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrui64List([]*uint64{})),
			value:   []any{uint64(3), uint64(2), uint64(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrui64List)
				for i, el := range list {
					expectedElem := expected[i].(uint64)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		// ptr int
		{
			name:    "int list marker to ptr int list",
			mrk:     parser.NewMarker("path:to:ptrintlist", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrintList([]*int{})),
			value:   []any{3, 2, 2},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrintList)
				for i, el := range list {
					expectedElem := expected[i].(int)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to ptr int8 list",
			mrk:     parser.NewMarker("path:to:ptri8list", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptri8List([]*int8{})),
			value:   []any{int8(3), int8(2), int8(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptri8List)
				for i, el := range list {
					expectedElem := expected[i].(int8)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to ptr int16 list",
			mrk:     parser.NewMarker("path:to:ptri16list", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptri16List([]*int16{})),
			value:   []any{int16(3), int16(2), int16(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptri16List)
				for i, el := range list {
					expectedElem := expected[i].(int16)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},

		{
			name:    "string list marker to ptr byte list",
			mrk:     parser.NewMarker("path:to:ptrbytelist", parser.MarkerKindList, reflect.ValueOf([]any{"p", "t", "m"})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrbyteList([]*byte{})),
			value:   []any{byte('p'), byte('t'), byte('m')},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrbyteList)
				for i, el := range list {
					expectedElem := expected[i].(byte)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to ptr byte list",
			mrk:     parser.NewMarker("path:to:ptrbytelist", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrbyteList([]*byte{})),
			value:   []any{byte(3), byte(2), byte(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrbyteList)
				for i, el := range list {
					expectedElem := expected[i].(byte)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "int list marker to ptr int64 list",
			mrk:     parser.NewMarker("path:to:ptri64list", parser.MarkerKindList, reflect.ValueOf([]any{3, 2, 2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptri64List([]*int64{})),
			value:   []any{int64(3), int64(2), int64(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptri64List)
				for i, el := range list {
					expectedElem := expected[i].(int64)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		// ptr float
		{
			name:    "float list marker to ptr float32 list",
			mrk:     parser.NewMarker("path:to:ptrf32list", parser.MarkerKindList, reflect.ValueOf([]any{3.0, 2.1, 2.2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrf32List([]*float32{})),
			value:   []any{float32(3.0), float32(2.1), float32(2.2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrf32List)
				for i, el := range list {
					expectedElem := expected[i].(float32)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "float list marker to ptr float64 list",
			mrk:     parser.NewMarker("path:to:ptrf64list", parser.MarkerKindList, reflect.ValueOf([]any{3.0, 2.1, 2.2})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrf64List([]*float64{})),
			value:   []any{float64(3.0), float64(2.1), float64(2.2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrf64List)
				for i, el := range list {
					expectedElem := expected[i].(float64)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		// complex
		{
			name:    "complex list marker to ptr complex64 list",
			mrk:     parser.NewMarker("path:to:ptrc64list", parser.MarkerKindList, reflect.ValueOf([]any{0 + 1i, 1 + 2i, 2 + 3i})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrc64List([]*complex64{})),
			value:   []any{complex64(0 + 1i), complex64(1 + 2i), complex64(2 + 3i)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrc64List)
				for i, el := range list {
					expectedElem := expected[i].(complex64)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "complex list marker to ptr complex128 list",
			mrk:     parser.NewMarker("path:to:ptrc128list", parser.MarkerKindList, reflect.ValueOf([]any{0 + 1i, 1 + 2i, 2 + 3i})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrc128List([]*complex128{})),
			value:   []any{0 + 1i, 1 + 2i, 2 + 3i},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrc128List)
				for i, el := range list {
					expectedElem := expected[i].(complex128)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		// bool
		{
			name:    "bool list marker to ptr bool list",
			mrk:     parser.NewMarker("path:to:ptrboollist", parser.MarkerKindList, reflect.ValueOf([]any{false, true, false})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrboolList([]*bool{})),
			value:   []any{false, true, false},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrboolList)
				for i, el := range list {
					expectedElem := expected[i].(bool)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
		// invalid test cases
		{
			name:    "string list marker to rune list wrong length",
			mrk:     parser.NewMarker("path:to:runelist", parser.MarkerKindList, reflect.ValueOf([]any{"p", "t", "ma"})),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(runeList([]rune{})),
			value:   nil,
			isValid: false,
			isValidValue: func(got reflect.Value, expected []any) bool {
				return expected != nil
			},
		},
	}
	reg := listDefs(t)
	mngr, err := NewConvMngr(reg)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, err := mngr.Convert(tc.mrk, tc.t)
			if err != nil && tc.isValid {
				t.Errorf("err occured: %s\n", err)
			}
			if err != nil && !tc.isValid {
				t.Skipf("wanted err occured: %v\n", err)
			}
			rtype := reflect.TypeOf(v)
			if rtype != tc.out {
				t.Fatalf("types don't match after conversion. got: %v; expected: %v\n", rtype, tc.out)
			}
			rvalue := reflect.ValueOf(v)
			if !tc.isValidValue(rvalue, tc.value) {
				t.Fatalf("value is not correct. got: %v", rvalue)
			}
			t.Logf("succesfully converted. got: %v; expected: %v\n", rtype, tc.out)
		})
	}
}
