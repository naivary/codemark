package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
)

// bytes
type byteList []byte

// rune
type runeList []rune

// string
type strList []string

// int
type intList []int
type i8List []int8
type i16List []int16
type i32List []int32
type i64List []int64

// uint
type uintList []uint
type ui8List []uint8
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
type ptri32List []*int32
type ptri64List []*int64

// ptr uint
type ptruintList []*uint
type ptrui8List []*uint8
type ptrui16List []*uint16
type ptrui32List []*uint32
type ptrui64List []*uint64

// float
type ptrf32List []*float32
type ptrf64List []*float64

// complex
type ptrc64List []*complex64
type ptrc128List []*complex128

// ptr bytes
type ptrbytesList []*byte

// ptr rune
type ptrruneList []*rune

func listDefs(t *testing.T) Registry {
	reg := NewInMemoryRegistry()
	defs := []*Definition{
		// list
		MakeDef("path:to:runelist", TargetField, reflect.TypeOf(runeList([]rune{}))),
		MakeDef("path:to:bytelist", TargetField, reflect.TypeOf(byteList([]byte{}))),
		MakeDef("path:to:stringlist", TargetField, reflect.TypeOf(strList([]string{}))),
		MakeDef("path:to:intlist", TargetField, reflect.TypeOf(intList([]int{}))),
		MakeDef("path:to:uintlist", TargetField, reflect.TypeOf(uintList([]uint{}))),
		MakeDef("path:to:f32list", TargetField, reflect.TypeOf(f32List([]float32{}))),
		MakeDef("path:to:f64list", TargetField, reflect.TypeOf(f64List([]float64{}))),
		MakeDef("path:to:c64list", TargetField, reflect.TypeOf(c64List([]complex64{}))),
		MakeDef("path:to:c128list", TargetField, reflect.TypeOf(c128List([]complex128{}))),
		MakeDef("path:to:boollist", TargetField, reflect.TypeOf(boolList([]bool{}))),
		// ptr
		MakeDef("path:to:ptrstringlist", TargetField, reflect.TypeOf(ptrstrList([]*string{}))),
		MakeDef("path:to:ptrintlist", TargetField, reflect.TypeOf(ptrintList([]*int{}))),
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
		t            Target
		out          reflect.Type
		value        []any
		isValid      bool
		isValidValue func(got reflect.Value, expected []any) bool
	}{
		{
			name:    "list marker to string list",
			mrk:     parser.NewMarker("path:to:stringlist", parser.MarkerKindList, reflect.ValueOf([]any{"path", "to", "marker"})),
			t:       TargetField,
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
		// pointer
		{
			name:    "list marker to ptr string list",
			mrk:     parser.NewMarker("path:to:ptrstringlist", parser.MarkerKindList, reflect.ValueOf([]any{"path", "to", "marker"})),
			t:       TargetField,
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
	}
	reg := listDefs(t)
	mngr, err := NewConvMngr(reg, &listConverter{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, err := mngr.Convert(tc.mrk, tc.t)
			if err != nil {
				t.Errorf("err occured: %s\n", err)
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
