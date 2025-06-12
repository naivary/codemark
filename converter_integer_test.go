package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

const _intValue = 99

func TestIntConverter(t *testing.T) {
	tests := []struct {
		name         string
		mrk          parser.Marker
		t            sdk.Target
		out          reflect.Type
		isValid      bool
		isValidValue func(got reflect.Value) bool
	}{
		// int
		{
			name:    "int marker to int type",
			mrk:     parser.NewMarker("path:to:i", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.Int(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.Int)
				return i == _intValue
			},
		},
		{
			name:    "int marker to int8 type",
			mrk:     parser.NewMarker("path:to:i8", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.I8(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.I8)
				return i == _intValue
			},
		},
		{
			name:    "int marker to int16 type",
			mrk:     parser.NewMarker("path:to:i16", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.I16(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.I16)
				return i == _intValue
			},
		},
		{
			name:    "int marker to int32 type",
			mrk:     parser.NewMarker("path:to:i32", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.I32(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.I32)
				return i == _intValue
			},
		},
		{
			name:    "int marker to int64 type",
			mrk:     parser.NewMarker("path:to:i64", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.I64(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.I64)
				return i == _intValue
			},
		},
		// uint
		{
			name:    "int marker to uint type",
			mrk:     parser.NewMarker("path:to:ui", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.Uint(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.Uint)
				return i == _intValue
			},
		},
		{
			name:    "int marker to uint8 type",
			mrk:     parser.NewMarker("path:to:ui8", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.U8(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.U8)
				return i == _intValue
			},
		},
		{
			name:    "int marker to uint16 type",
			mrk:     parser.NewMarker("path:to:ui16", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.U16(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.U16)
				return i == _intValue
			},
		},
		{
			name:    "int marker to uint32 type",
			mrk:     parser.NewMarker("path:to:ui32", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.U32(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.U32)
				return i == _intValue
			},
		},
		{
			name:    "int marker to uint64 type",
			mrk:     parser.NewMarker("path:to:ui64", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.U64(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.U64)
				return i == _intValue
			},
		},
		// ptr int
		{
			name:    "int marker to ptr int type",
			mrk:     parser.NewMarker("path:to:ptri", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrInt(new(int))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.PtrInt)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr int8 type",
			mrk:     parser.NewMarker("path:to:ptri8", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrI8(new(int8))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.PtrI8)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr int16 type",
			mrk:     parser.NewMarker("path:to:ptri16", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrI16(new(int16))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.PtrI16)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr int32 type",
			mrk:     parser.NewMarker("path:to:ptri32", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrI32(new(int32))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.PtrI32)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr int64 type",
			mrk:     parser.NewMarker("path:to:ptri64", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrI64(new(int64))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.PtrI64)
				return *i == _intValue
			},
		},
		// ptr uint
		{
			name:    "int marker to ptr uint type",
			mrk:     parser.NewMarker("path:to:ptrui", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrUint(new(uint))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.PtrUint)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr uint8 type",
			mrk:     parser.NewMarker("path:to:ptrui8", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrU8(new(uint8))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.PtrU8)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr uint16 type",
			mrk:     parser.NewMarker("path:to:ptrui16", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrU16(new(uint16))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.PtrU16)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr uint32 type",
			mrk:     parser.NewMarker("path:to:ptrui32", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrU32(new(uint32))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.PtrU32)
				return *i == _intValue
			},
		},
		{
			name:    "int marker to ptr uint64 type",
			mrk:     parser.NewMarker("path:to:ptrui64", parser.MarkerKindInt, reflect.ValueOf(_intValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrU64(new(uint64))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				i := got.Interface().(sdktesting.PtrU64)
				return *i == _intValue
			},
		},
		{
			name:    "string marker to byte type",
			mrk:     parser.NewMarker("path:to:byte", parser.MarkerKindString, reflect.ValueOf("c")),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.Byte(0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				b := got.Interface().(sdktesting.Byte)
				return b == 'c'
			},
		},
	}
	reg, err := sdktesting.NewDefsSet(NewInMemoryRegistry(), &DefinitionMarker{})
	if err != nil {
		t.Errorf("err occured: %s\n", err.Error())
	}
	mngr, err := NewConvMngr(reg)
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
