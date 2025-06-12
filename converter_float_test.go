package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func TestFloatConverter(t *testing.T) {
	tests := []sdktesting.ConverterTestCase{
		{
			Name:        "float marker to f32 type",
			Marker:      sdktesting.RandMarker(sdktesting.F32(0.0)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.F32(0.0)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := float64(got.Interface().(sdktesting.F32))
				w := wanted.Interface().(float64)
				return sdktesting.AlmostEqual(g, w)
			},
		},
		{
			Name:        "float marker to f64 type",
			Marker:      sdktesting.RandMarker(sdktesting.F64(0.0)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.F64(0.0)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := got.Interface().(sdktesting.F64)
				w := wanted.Interface().(float64)
				return sdktesting.AlmostEqual(float64(g), w)
			},
		},
		{
			Name:        "float marker to ptr f32 type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrF32(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrF32(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := float64(*got.Interface().(sdktesting.PtrF32))
				w := wanted.Interface().(float64)
				return sdktesting.AlmostEqual(g, w)
			},
		},
		{
			Name:        "float marker to ptr f64 type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrF64(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrF64(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := float64(*got.Interface().(sdktesting.PtrF64))
				w := wanted.Interface().(float64)
				return sdktesting.AlmostEqual(g, w)
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
		t.Run(tc.Name, func(t *testing.T) {
			v, err := mngr.Convert(tc.Marker, tc.Target)
			if err != nil && tc.IsValidCase {
				t.Errorf("err occured: %s\n", err)
			}
			if err != nil && !tc.IsValidCase {
				t.Skipf("skipping invalid case after successful assertion")
			}
			gotType := reflect.TypeOf(v)
			if gotType != tc.ToType {
				t.Fatalf("types don't match after conversion. got: %v; expected: %v\n", gotType, tc.ToType)
			}
			gotValue := reflect.ValueOf(v)
			if !tc.IsValidValue(gotValue, tc.Marker.Value()) {
				t.Fatalf("value is not correct. got: %v; wanted: %v\n", gotValue, tc.Marker.Value())
			}
			t.Logf("succesfully converted. got: %v; expected: %v\n", gotType, tc.ToType)
		})
	}
}
