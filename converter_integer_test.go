package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
)

const _intValue = 99

type i int
type i8 int8
type i16 int16
type i32 int32
type i64 int64

type ptri *int
type ptri8 *int8
type ptri16 *int16
type ptri32 *int32
type ptri64 *int64

type ui uint
type ui8 uint8
type ui16 uint16
type ui32 uint32
type ui64 uint64

type ptrui *uint
type ptrui8 *uint8
type ptrui16 *uint16
type ptrui32 *uint32
type ptrui64 *uint64

func intDefs(t *testing.T) Registry {
	reg := NewInMemoryRegistry()
	defs := []*Definition{
		// int
		MakeDef("path:to:i", TargetField, reflect.TypeOf(i(0))),
		MakeDef("path:to:i8", TargetField, reflect.TypeOf(i8(0))),
		MakeDef("path:to:i16", TargetField, reflect.TypeOf(i16(0))),
		MakeDef("path:to:i32", TargetField, reflect.TypeOf(i32(0))),
		MakeDef("path:to:i64", TargetField, reflect.TypeOf(i64(0))),
		// ptr int
		MakeDef("path:to:ptri", TargetField, reflect.TypeOf(ptri(new(int)))),
		MakeDef("path:to:ptri8", TargetField, reflect.TypeOf(ptri8(new(int8)))),
		MakeDef("path:to:ptri16", TargetField, reflect.TypeOf(ptri16(new(int16)))),
		MakeDef("path:to:ptri32", TargetField, reflect.TypeOf(ptri32(new(int32)))),
		MakeDef("path:to:ptri64", TargetField, reflect.TypeOf(ptri64(new(int64)))),
		// uint
		MakeDef("path:to:ui", TargetField, reflect.TypeOf(ui(0))),
		MakeDef("path:to:ui8", TargetField, reflect.TypeOf(ui8(0))),
		MakeDef("path:to:ui16", TargetField, reflect.TypeOf(ui16(0))),
		MakeDef("path:to:ui32", TargetField, reflect.TypeOf(ui32(0))),
		MakeDef("path:to:ui64", TargetField, reflect.TypeOf(ui64(0))),
		// ptr uint
		MakeDef("path:to:ptrui", TargetField, reflect.TypeOf(ptrui(new(uint)))),
		MakeDef("path:to:ptrui8", TargetField, reflect.TypeOf(ptrui8(new(uint8)))),
		MakeDef("path:to:ptrui16", TargetField, reflect.TypeOf(ptrui16(new(uint16)))),
		MakeDef("path:to:ptrui32", TargetField, reflect.TypeOf(ptrui32(new(uint32)))),
		MakeDef("path:to:ptrui64", TargetField, reflect.TypeOf(ptrui64(new(uint64)))),
	}
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
			t.Errorf("err occured: %s\n", err)
		}

	}
	return reg
}

func TestIntConverter(t *testing.T) {
	tests := []struct {
		name         string
		mrk          parser.Marker
		t            Target
		out          reflect.Type
		isValid      bool
		isValidValue func(got reflect.Value) bool
	}{
		// int
		{
			name:    "int marker to int type",
			mrk:     parser.NewMarker("path:to:i", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(i(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(i)
				return i == _intValue
			},
		},
		{
			name:    "int marker to int8 type",
			mrk:     parser.NewMarker("path:to:i8", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(i8(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(i8)
				return i == _intValue
			},
		},
		{
			name:    "int marker to int16 type",
			mrk:     parser.NewMarker("path:to:i16", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(i16(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(i16)
				return i == _intValue
			},
		},
		{
			name:    "int marker to int32 type",
			mrk:     parser.NewMarker("path:to:i32", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(i32(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(i32)
				return i == _intValue
			},
		},
		{
			name:    "int marker to int64 type",
			mrk:     parser.NewMarker("path:to:i64", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(i64(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(i64)
				return i == _intValue
			},
		},
		// uint
		{
			name:    "int marker to uint type",
			mrk:     parser.NewMarker("path:to:ui", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ui(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ui)
				return i == _intValue
			},
		},
		{
			name:    "int marker to uint8 type",
			mrk:     parser.NewMarker("path:to:ui8", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ui8(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ui8)
				return i == _intValue
			},
		},
		{
			name:    "int marker to uint16 type",
			mrk:     parser.NewMarker("path:to:ui16", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ui16(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ui16)
				return i == _intValue
			},
		},
		{
			name:    "int marker to uint32 type",
			mrk:     parser.NewMarker("path:to:ui32", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ui32(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ui32)
				return i == _intValue
			},
		},
		{
			name:    "int marker to uint64 type",
			mrk:     parser.NewMarker("path:to:ui64", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ui64(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ui64)
				return i == _intValue
			},
		},
		// ptr int
		{
			name:    "int marker to ptr int type",
			mrk:     parser.NewMarker("path:to:ptri", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptri(new(int))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ptri)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr int8 type",
			mrk:     parser.NewMarker("path:to:ptri8", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptri8(new(int8))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ptri8)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr int16 type",
			mrk:     parser.NewMarker("path:to:ptri16", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptri16(new(int16))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ptri16)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr int32 type",
			mrk:     parser.NewMarker("path:to:ptri32", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptri32(new(int32))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ptri32)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr int64 type",
			mrk:     parser.NewMarker("path:to:ptri64", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptri64(new(int64))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ptri64)
				return *i == _intValue
			},
		},
		// ptr uint
		{
			name:    "int marker to ptr uint type",
			mrk:     parser.NewMarker("path:to:ptrui", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptrui(new(uint))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ptrui)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr uint8 type",
			mrk:     parser.NewMarker("path:to:ptrui8", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptrui8(new(uint8))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ptrui8)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr uint16 type",
			mrk:     parser.NewMarker("path:to:ptrui16", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptrui16(new(uint16))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ptrui16)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr uint32 type",
			mrk:     parser.NewMarker("path:to:ptrui32", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptrui32(new(uint32))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ptrui32)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr uint64 type",
			mrk:     parser.NewMarker("path:to:ptrui64", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptrui64(new(uint64))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(ptrui64)
				return *i == _intValue
			},
		},
	}
	reg := intDefs(t)
	mngr, err := NewConvMngr(reg, &intConverter{})
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
			if !tc.isValidValue(rvalue) {
				t.Fatalf("value is not correct. got: %v", rvalue)
			}
			t.Logf("succesfully converted. got: %v; expected: %v\n", rtype, tc.out)
		})
	}
}
