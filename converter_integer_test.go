package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

const _intValue = 99

func TestIntConverter(t *testing.T) {
	tests := []sdktesting.ConverterTestCase{
		// int
		{
			Name:        "int marker to int type",
			Marker:      sdktesting.RandMarker(sdktesting.Int(0)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.Int(0)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := int64(got.Interface().(sdktesting.Int))
				w := wanted.Interface().(int64)
				return g == w
			},
		},
		{
			Name:        "int marker to int8 type",
			Marker:      sdktesting.RandMarker(sdktesting.I8(0)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.I8(0)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := int64(got.Interface().(sdktesting.I8))
				w := wanted.Interface().(int64)
				return g == w
			},
		},
		{
			Name:        "int marker to int16 type",
			Marker:      sdktesting.RandMarker(sdktesting.I16(0)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.I16(0)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := int64(got.Interface().(sdktesting.I16))
				w := wanted.Interface().(int64)
				return g == w
			},
		},
		{
			Name:        "int marker to int32 type",
			Marker:      sdktesting.RandMarker(sdktesting.I32(0)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.I32(0)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := int64(got.Interface().(sdktesting.I32))
				w := wanted.Interface().(int64)
				return g == w
			},
		},
		{
			Name:        "int marker to int64 type",
			Marker:      sdktesting.RandMarker(sdktesting.I64(0)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.I64(0)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := int64(got.Interface().(sdktesting.I64))
				w := wanted.Interface().(int64)
				return g == w
			},
		},
		// uint
		{
			Name:        "int marker to uint type",
			Marker:      sdktesting.RandMarker(sdktesting.Uint(0)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.Uint(0)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := uint64(got.Interface().(sdktesting.Uint))
				w := wanted.Interface().(int64)
				return g == uint64(w)
			},
		},
		{
			Name:        "int marker to uint8 type",
			Marker:      sdktesting.RandMarker(sdktesting.U8(0)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.U8(0)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := uint64(got.Interface().(sdktesting.U8))
				w := wanted.Interface().(int64)
				return g == uint64(w)
			},
		},
		{
			Name:        "int marker to uint16 type",
			Marker:      sdktesting.RandMarker(sdktesting.U16(0)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.U16(0)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := uint64(got.Interface().(sdktesting.U16))
				w := wanted.Interface().(int64)
				return g == uint64(w)
			},
		},
		{
			Name:        "int marker to uint32 type",
			Marker:      sdktesting.RandMarker(sdktesting.U32(0)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.U32(0)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := uint64(got.Interface().(sdktesting.U32))
				w := wanted.Interface().(int64)
				return g == uint64(w)
			},
		},
		{
			Name:        "int marker to uint64 type",
			Marker:      sdktesting.RandMarker(sdktesting.U64(0)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.U64(0)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := uint64(got.Interface().(sdktesting.U64))
				w := wanted.Interface().(int64)
				return g == uint64(w)
			},
		},
		// ptr int
		{
			Name:        "int marker to ptr int type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrInt(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrInt(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := int64(*got.Interface().(sdktesting.PtrInt))
				w := wanted.Interface().(int64)
				return g == w
			},
		},
		{
			Name:        "int marker to ptr int8 type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrI8(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrI8(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := int64(*got.Interface().(sdktesting.PtrI8))
				w := wanted.Interface().(int64)
				return g == w
			},
		},
		{
			Name:        "int marker to ptr int16 type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrI16(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrI16(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := int64(*got.Interface().(sdktesting.PtrI16))
				w := wanted.Interface().(int64)
				return g == w
			},
		},
		{
			Name:        "int marker to ptr int32 type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrI32(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrI32(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := int64(*got.Interface().(sdktesting.PtrI32))
				w := wanted.Interface().(int64)
				return g == w
			},
		},
		{
			Name:        "int marker to ptr int64 type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrI64(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrI64(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := int64(*got.Interface().(sdktesting.PtrI64))
				w := wanted.Interface().(int64)
				return g == w
			},
		},
		// ptr uint
		{
			Name:        "int marker to ptr uint type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrUint(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrUint(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := uint64(*got.Interface().(sdktesting.PtrUint))
				w := wanted.Interface().(int64)
				return g == uint64(w)
			},
		},
		{
			Name:        "int marker to ptr uint8 type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrU8(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrU8(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := uint64(*got.Interface().(sdktesting.PtrU8))
				w := wanted.Interface().(int64)
				return g == uint64(w)
			},
		},
		{
			Name:        "int marker to ptr uint16 type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrU16(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrU16(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := uint64(*got.Interface().(sdktesting.PtrU16))
				w := wanted.Interface().(int64)
				return g == uint64(w)
			},
		},
		{
			Name:        "int marker to ptr uint32 type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrU32(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrU32(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := uint64(*got.Interface().(sdktesting.PtrU32))
				w := wanted.Interface().(int64)
				return g == uint64(w)
			},
		},
		{
			Name:        "int marker to ptr uint64 type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrU64(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrU64(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := uint64(*got.Interface().(sdktesting.PtrU64))
				w := wanted.Interface().(int64)
				return g == uint64(w)
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
		t.Run(tc.Name, func(t *testing.T) {
			v, err := mngr.Convert(tc.Marker, tc.Target)
			if err != nil {
				t.Errorf("err occured: %s\n", err)
			}
			gotType := reflect.TypeOf(v)
			if gotType != tc.ToType {
				t.Fatalf("types don't match after conversion. got: %v; expected: %v\n", gotType, tc.ToType)
			}
			gotValue := reflect.ValueOf(v)
			if !tc.IsValidValue(gotValue, tc.Marker.Value()) {
				t.Fatalf("value is not correct. got: %v; expected: %v\n", gotValue, tc.Marker.Value())
			}
			t.Logf("succesfully converted. got: %v; expected: %v\n", gotType, tc.ToType)
		})
	}
}
