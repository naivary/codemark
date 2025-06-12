package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

const _floatValue = 99.99

func TestFloatConverter(t *testing.T) {
	tests := []struct {
		name         string
		mrk          parser.Marker
		t            sdk.Target
		out          reflect.Type
		isValid      bool
		isValidValue func(got reflect.Value) bool
	}{
		{
			name:    "float marker to f32 type",
			mrk:     parser.NewMarker("path:to:f32", parser.MarkerKindFloat, reflect.ValueOf(_floatValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.F32(0.0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				f := got.Interface().(sdktesting.F32)
				return f == _floatValue
			},
		},
		{
			name:    "float marker to f64 type",
			mrk:     parser.NewMarker("path:to:f64", parser.MarkerKindFloat, reflect.ValueOf(_floatValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.F64(0.0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				f := got.Interface().(sdktesting.F64)
				return f == _floatValue
			},
		},
		{
			name:    "float marker to ptr f32 type",
			mrk:     parser.NewMarker("path:to:ptrf32", parser.MarkerKindFloat, reflect.ValueOf(_floatValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrF32(new(float32))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				f := got.Interface().(sdktesting.PtrF32)
				return *f == _floatValue
			},
		},
		{
			name:    "float marker to ptr f64 type",
			mrk:     parser.NewMarker("path:to:ptrf64", parser.MarkerKindFloat, reflect.ValueOf(_floatValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrF64(new(float64))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				f := got.Interface().(sdktesting.PtrF64)
				return *f == _floatValue
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
