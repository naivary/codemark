package codemark

import (
	"math"
	"reflect"
	"testing"

	"github.com/naivary/codemark/marker"
)

func ptr[T any](t T) *T {
	return &t
}

type str string
type strPtr *string
type stringSlice []string
type stringSliceElemPtr []*string
type stringSliceDoublePtr *[]*string
type stringSlicePtr *[]string

type c64 complex64
type c128 complex128
type c64Slice []complex64
type c128Slice []complex128

type f32 float32
type f64 float64
type f32Slice []float32
type f64Slice []float64

type i int
type i8 int8
type i16 int16
type i32 int32
type i64 int64

type iSlice []int
type iSliceElemPtr []*int
type iSliceDoublePtr *[]*int
type iSlicePtr *[]int

type i8Slice []int8
type i16Slice []int16
type i32Slice []int32
type i64Slice []int64

type iPtr *int
type i8Ptr *int8
type i16Ptr *int16
type i32Ptr *int32
type i64Ptr *int64

type ui uint
type ui8 uint8
type ui16 uint16
type ui32 uint32
type ui64 uint64

type b byte
type bytePtr *byte
type r rune
type runePtr *rune

type byteSlice []byte
type byteSlicePtr *[]byte

type runeSlice []rune
type runeSlicePtr *[]rune

type boolean bool
type boolPtr *bool

type anySlice []any

func TestConverter_Convert(t *testing.T) {
	tests := []struct {
		name    string
		markers []marker.Marker
		defs    []*Definition
	}{
		{
			name: "string marker",
			markers: []marker.Marker{
				marker.NewDefault("jsonschema:validation:string", reflect.String, reflect.ValueOf("string")),
				marker.NewDefault("jsonschema:validation:byte", reflect.String, reflect.ValueOf("b")),
				marker.NewDefault("jsonschema:validation:rune", reflect.String, reflect.ValueOf("r")),
			},
			defs: []*Definition{
				MakeDef("jsonschema:validation:string", TargetConst, reflect.TypeOf(str(""))),
				MakeDef("jsonschema:validation:byte", TargetConst, reflect.TypeOf(b(0))),
				MakeDef("jsonschema:validation:rune", TargetConst, reflect.TypeOf(r(0))),
			},
		},
		{
			name: "boolean marker",
			markers: []marker.Marker{
				marker.NewDefault("jsonschema:validation:required", reflect.Bool, reflect.ValueOf(true)),
				marker.NewDefault("jsonschema:validation:requiredptr", reflect.Bool, reflect.ValueOf(true)),
			},
			defs: []*Definition{
				MakeDef("jsonschema:validation:required", TargetConst, reflect.TypeOf(boolean(false))),
				MakeDef("jsonschema:validation:requiredptr", TargetConst, reflect.TypeOf(boolPtr(ptr(false)))),
			},
		},
		{
			name: "string marker pointers",
			markers: []marker.Marker{
				marker.NewDefault("jsonschema:validation:stringptr", reflect.String, reflect.ValueOf("stringptr")),
				marker.NewDefault("jsonschema:validation:byteptr", reflect.String, reflect.ValueOf("b")),
				marker.NewDefault("jsonschema:validation:runeptr", reflect.String, reflect.ValueOf("r")),
			},
			defs: []*Definition{
				MakeDef("jsonschema:validation:stringptr", TargetConst, reflect.TypeOf(strPtr(ptr("")))),
				MakeDef("jsonschema:validation:byteptr", TargetConst, reflect.TypeOf(bytePtr(ptr[byte](0)))),
				MakeDef("jsonschema:validation:runeptr", TargetConst, reflect.TypeOf(runePtr(ptr[rune](0)))),
			},
		},
		{
			name: "integer marker",
			defs: []*Definition{
				MakeDef("jsonschema:validation:i", TargetConst, reflect.TypeOf(i(0))),
				MakeDef("jsonschema:validation:i8", TargetConst, reflect.TypeOf(i8(0))),
				MakeDef("jsonschema:validation:i16", TargetConst, reflect.TypeOf(i16(0))),
				MakeDef("jsonschema:validation:i32", TargetConst, reflect.TypeOf(i32(0))),
				MakeDef("jsonschema:validation:i64", TargetConst, reflect.TypeOf(i64(0))),
				MakeDef("jsonschema:validation:rune", TargetConst, reflect.TypeOf(r(0))),
			},
			markers: []marker.Marker{
				marker.NewDefault("jsonschema:validation:i", reflect.Int64, reflect.ValueOf(30)),
				marker.NewDefault("jsonschema:validation:i8", reflect.Int64, reflect.ValueOf(math.MaxInt8)),
				marker.NewDefault("jsonschema:validation:i8", reflect.Int64, reflect.ValueOf(2)),
				marker.NewDefault("jsonschema:validation:i16", reflect.Int64, reflect.ValueOf(32767)),
				marker.NewDefault("jsonschema:validation:i32", reflect.Int64, reflect.ValueOf(213712389)),
				marker.NewDefault("jsonschema:validation:i64", reflect.Int64, reflect.ValueOf(1237812378192378)),
				marker.NewDefault("jsonschema:validation:rune", reflect.Int64, reflect.ValueOf(math.MaxInt32)),
			},
		},
		{
			name: "float marker",
			markers: []marker.Marker{
				marker.NewDefault("jsonschema:validation:f32", reflect.Float64, reflect.ValueOf(30.32)),
				marker.NewDefault("jsonschema:validation:f64", reflect.Float64, reflect.ValueOf(2.213)),
				marker.NewDefault("jsonschema:validation:f32", reflect.Float64, reflect.ValueOf(-30.32)),
			},
			defs: []*Definition{
				MakeDef("jsonschema:validation:f32", TargetConst, reflect.TypeOf(f32(0.0))),
				MakeDef("jsonschema:validation:f64", TargetConst, reflect.TypeOf(f64(0.0))),
			},
		},
		{
			name: "complex marker",
			markers: []marker.Marker{
				marker.NewDefault("jsonschema:validation:c64", reflect.Complex128, reflect.ValueOf(3+2i)),
				marker.NewDefault("jsonschema:validation:c128", reflect.Complex128, reflect.ValueOf(3+3i)),
			},
			defs: []*Definition{
				MakeDef("jsonschema:validation:c64", TargetConst, reflect.TypeOf(c64(1+1i))),
				MakeDef("jsonschema:validation:c128", TargetConst, reflect.TypeOf(c128(1+1i))),
			},
		},
		{
			name: "slice marker",
			defs: []*Definition{
				MakeDef("jsonschema:validation:stringselemptr", TargetConst, reflect.TypeOf(stringSliceElemPtr([]*string{}))),
				MakeDef("jsonschema:validation:stringsdoubleptr", TargetConst, reflect.TypeOf(stringSliceDoublePtr(new([]*string)))),
				MakeDef("jsonschema:validation:stringsptr", TargetConst, reflect.TypeOf(stringSlicePtr(new([]string)))),
				MakeDef("jsonschema:validation:stringslice", TargetConst, reflect.TypeOf(stringSlice([]string{}))),
				MakeDef("jsonschema:validation:ints", TargetConst, reflect.TypeOf(iSlice([]int{}))),
				MakeDef("jsonschema:validation:ints64", TargetConst, reflect.TypeOf(i64Slice([]int64{}))),
				MakeDef("jsonschema:validation:intselemptr", TargetConst, reflect.TypeOf(iSliceElemPtr([]*int{}))),
				MakeDef("jsonschema:validation:intdoubleptr", TargetConst, reflect.TypeOf(iSliceDoublePtr(new([]*int)))),
				MakeDef("jsonschema:validation:intsptr", TargetConst, reflect.TypeOf(iSlicePtr(new([]int)))),
				MakeDef("jsonschema:validation:bytes", TargetConst, reflect.TypeOf(byteSlice([]byte{}))),
				MakeDef("jsonschema:validation:anys", TargetConst, reflect.TypeOf(anySlice([]any{}))),
			},
			markers: []marker.Marker{
				marker.NewDefault("jsonschema:validation:stringsdoubleptr", reflect.Slice, reflect.ValueOf([]any{"first", "onemore", "again"})),
				marker.NewDefault("jsonschema:validation:stringselemptr", reflect.Slice, reflect.ValueOf([]any{"first", "onemore", "again"})),
				marker.NewDefault("jsonschema:validation:stringsptr", reflect.Slice, reflect.ValueOf([]any{"first", "onemore", "again"})),
				marker.NewDefault("jsonschema:validation:stringslice", reflect.Slice, reflect.ValueOf([]any{"without", "pointer", "slice"})),
				marker.NewDefault("jsonschema:validation:ints", reflect.Slice, reflect.ValueOf([]any{2, 3, 0x321})),
				marker.NewDefault("jsonschema:validation:ints64", reflect.Slice, reflect.ValueOf([]any{2, 3, 0x321})),
				marker.NewDefault("jsonschema:validation:intselemptr", reflect.Slice, reflect.ValueOf([]any{2, 3, 0x321})),
				marker.NewDefault("jsonschema:validation:intdoubleptr", reflect.Slice, reflect.ValueOf([]any{2, 3, 0x321})),
				marker.NewDefault("jsonschema:validation:intsptr", reflect.Slice, reflect.ValueOf([]any{2, 3, 0x321})),
				marker.NewDefault("jsonschema:validation:bytes", reflect.Slice, reflect.ValueOf([]any{"s", "d", "a"})),
				marker.NewDefault("jsonschema:validation:anys", reflect.Slice, reflect.ValueOf([]any{86, 3, 2})),
				marker.NewDefault("jsonschema:validation:anys", reflect.Slice, reflect.ValueOf([]any{86, "str", 2})),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reg := NewRegistry()
			for _, def := range tc.defs {
				if err := reg.Define(def); err != nil {
					t.Error(err)
				}
			}
			conv, err := NewConverter(reg)
			if err != nil {
				t.Error(err)
			}
			for _, marker := range tc.markers {
				v, err := conv.Convert(marker, TargetConst)
				if err != nil {
					t.Fatal(err)
				}
				t.Log(v)
				def := reg.Get(marker.Ident())
				typ := reflect.TypeOf(v)
				if typ != def.output {
					t.Fatalf("conversion to definition type did not work. Expected: %v. Got: %v", def.output, typ)
				}
			}

		})
	}

}
