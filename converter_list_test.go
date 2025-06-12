package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

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
			out:     reflect.TypeOf(sdktesting.StringList([]string{})),
			value:   []any{"path", "to", "marker"},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.StringList)
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
			out:     reflect.TypeOf(sdktesting.RuneList([]rune{})),
			value:   []any{rune('p'), rune('t'), rune('m')},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.RuneList)
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
			out:     reflect.TypeOf(sdktesting.UintList([]uint{})),
			value:   []any{uint(3), uint(2), uint(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.UintList)
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
			out:     reflect.TypeOf(sdktesting.RuneList([]rune{})),
			value:   []any{rune(3), rune(2), rune(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.RuneList)
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
			out:     reflect.TypeOf(sdktesting.U16List([]uint16{})),
			value:   []any{uint16(3), uint16(2), uint16(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.U16List)
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
			out:     reflect.TypeOf(sdktesting.U32List([]uint32{})),
			value:   []any{uint32(3), uint32(2), uint32(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.U32List)
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
			out:     reflect.TypeOf(sdktesting.U64List([]uint64{})),
			value:   []any{uint64(3), uint64(2), uint64(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.U64List)
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
			out:     reflect.TypeOf(sdktesting.IntList([]int{})),
			value:   []any{3, 2, 2},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.IntList)
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
			out:     reflect.TypeOf(sdktesting.I8List([]int8{})),
			value:   []any{int8(3), int8(2), int8(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.I8List)
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
			out:     reflect.TypeOf(sdktesting.I16List([]int16{})),
			value:   []any{int16(3), int16(2), int16(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.I16List)
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
			out:     reflect.TypeOf(sdktesting.ByteList([]byte{})),
			value:   []any{byte('p'), byte('t'), byte('m')},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.ByteList)
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
			out:     reflect.TypeOf(sdktesting.ByteList([]byte{})),
			value:   []any{byte(3), byte(2), byte(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.ByteList)
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
			out:     reflect.TypeOf(sdktesting.I64List([]int64{})),
			value:   []any{int64(3), int64(2), int64(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.I64List)
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
			out:     reflect.TypeOf(sdktesting.F32List([]float32{})),
			value:   []any{float32(3.0), float32(2.1), float32(2.2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.F32List)
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
			out:     reflect.TypeOf(sdktesting.F64List([]float64{})),
			value:   []any{float64(3.0), float64(2.1), float64(2.2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.F64List)
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
			out:     reflect.TypeOf(sdktesting.C64List([]complex64{})),
			value:   []any{complex64(0 + 1i), complex64(1 + 2i), complex64(2 + 3i)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.C64List)
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
			out:     reflect.TypeOf(sdktesting.C128List([]complex128{})),
			value:   []any{0 + 1i, 1 + 2i, 2 + 3i},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.C128List)
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
			out:     reflect.TypeOf(sdktesting.BoolList([]bool{})),
			value:   []any{false, true, false},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.BoolList)
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
			out:     reflect.TypeOf(sdktesting.PtrStringList([]*string{})),
			value:   []any{"path", "to", "marker"},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrStringList)
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
			out:     reflect.TypeOf(sdktesting.PtrUintList([]*uint{})),
			value:   []any{uint(3), uint(2), uint(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrUintList)
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
			out:     reflect.TypeOf(sdktesting.PtrRuneList([]*rune{})),
			value:   []any{rune('p'), rune('t'), rune('m')},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrRuneList)
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
			out:     reflect.TypeOf(sdktesting.PtrRuneList([]*rune{})),
			value:   []any{rune(3), rune(2), rune(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrRuneList)
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
			out:     reflect.TypeOf(sdktesting.PtrU16List([]*uint16{})),
			value:   []any{uint16(3), uint16(2), uint16(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrU16List)
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
			out:     reflect.TypeOf(sdktesting.PtrU32List([]*uint32{})),
			value:   []any{uint32(3), uint32(2), uint32(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrU32List)
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
			out:     reflect.TypeOf(sdktesting.PtrU64List([]*uint64{})),
			value:   []any{uint64(3), uint64(2), uint64(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrU64List)
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
			out:     reflect.TypeOf(sdktesting.PtrIntList([]*int{})),
			value:   []any{3, 2, 2},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrIntList)
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
			out:     reflect.TypeOf(sdktesting.PtrI8List([]*int8{})),
			value:   []any{int8(3), int8(2), int8(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrI8List)
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
			out:     reflect.TypeOf(sdktesting.PtrI16List([]*int16{})),
			value:   []any{int16(3), int16(2), int16(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrI16List)
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
			out:     reflect.TypeOf(sdktesting.PtrByteList([]*byte{})),
			value:   []any{byte('p'), byte('t'), byte('m')},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrByteList)
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
			out:     reflect.TypeOf(sdktesting.PtrByteList([]*byte{})),
			value:   []any{byte(3), byte(2), byte(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrByteList)
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
			out:     reflect.TypeOf(sdktesting.PtrI64List([]*int64{})),
			value:   []any{int64(3), int64(2), int64(2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrI64List)
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
			out:     reflect.TypeOf(sdktesting.PtrF32List([]*float32{})),
			value:   []any{float32(3.0), float32(2.1), float32(2.2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrF32List)
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
			out:     reflect.TypeOf(sdktesting.PtrF64List([]*float64{})),
			value:   []any{float64(3.0), float64(2.1), float64(2.2)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrF64List)
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
			out:     reflect.TypeOf(sdktesting.PtrC64List([]*complex64{})),
			value:   []any{complex64(0 + 1i), complex64(1 + 2i), complex64(2 + 3i)},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrC64List)
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
			out:     reflect.TypeOf(sdktesting.PtrC128List([]*complex128{})),
			value:   []any{0 + 1i, 1 + 2i, 2 + 3i},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrC128List)
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
			out:     reflect.TypeOf(sdktesting.PtrBoolList([]*bool{})),
			value:   []any{false, true, false},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(sdktesting.PtrBoolList)
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
			out:     reflect.TypeOf(sdktesting.RuneList([]rune{})),
			value:   nil,
			isValid: false,
			isValidValue: func(got reflect.Value, expected []any) bool {
				return expected != nil
			},
		},
	}
	reg, err := sdktesting.NewDefsSet(NewInMemoryRegistry(), &DefinitionMarker{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
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
