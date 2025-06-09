package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
)

const _floatValue = 99.99

type f32 float32
type f64 float64

type ptrf32 *float32
type ptrf64 *float64

func floatDefs(t *testing.T) Registry {
	reg := NewInMemoryRegistry()
	defs := []*Definition{
		// float
		MakeDef("path:to:f32", TargetField, reflect.TypeOf(f32(0.0))),
		MakeDef("path:to:f64", TargetField, reflect.TypeOf(f64(0.0))),
		// ptr float
		MakeDef("path:to:ptrf32", TargetField, reflect.TypeOf(ptrf32(new(float32)))),
		MakeDef("path:to:ptrf64", TargetField, reflect.TypeOf(ptrf64(new(float64)))),
	}
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
			t.Errorf("err occured: %s\n", err)
		}

	}
	return reg
}

func TestFloatConverter(t *testing.T) {
	tests := []struct {
		name         string
		mrk          parser.Marker
		t            Target
		out          reflect.Type
		isValid      bool
		isValidValue func(got reflect.Value) bool
	}{
		{
			name:    "float marker to f32 type",
			mrk:     parser.NewMarker("path:to:f32", parser.MarkerKindFloat, reflect.ValueOf(_floatValue)),
			t:       TargetField,
			out:     reflect.TypeOf(f32(0.0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				f := got.Interface().(f32)
				return f == _floatValue
			},
		},
		{
			name:    "float marker to f64 type",
			mrk:     parser.NewMarker("path:to:f64", parser.MarkerKindFloat, reflect.ValueOf(_floatValue)),
			t:       TargetField,
			out:     reflect.TypeOf(f64(0.0)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				f := got.Interface().(f64)
				return f == _floatValue
			},
		},
		{
			name:    "float marker to ptr f32 type",
			mrk:     parser.NewMarker("path:to:ptrf32", parser.MarkerKindFloat, reflect.ValueOf(_floatValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptrf32(new(float32))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				f := got.Interface().(ptrf32)
				return *f == _floatValue
			},
		},
		{
			name:    "float marker to ptr f64 type",
			mrk:     parser.NewMarker("path:to:ptrf64", parser.MarkerKindFloat, reflect.ValueOf(_floatValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptrf64(new(float64))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				f := got.Interface().(ptrf64)
				return *f == _floatValue
			},
		},

	}
	reg := floatDefs(t)
	mngr, err := NewConvMngr(reg, &floatConverter{})
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
