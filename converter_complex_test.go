package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

const _complexValue = 9 + 9i

func TestComplexConverter(t *testing.T) {
	tests := []struct {
		name         string
		mrk          parser.Marker
		t            sdk.Target
		out          reflect.Type
		isValid      bool
		isValidValue func(got reflect.Value) bool
	}{
		{
			name:    "complex marker to c64 type",
			mrk:     parser.NewMarker("path:to:c64", parser.MarkerKindComplex, reflect.ValueOf(_complexValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.C64(0 + 0i)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				c := got.Interface().(sdktesting.C64)
				return c == _complexValue
			},
		},
		{
			name:    "complex marker to c128 type",
			mrk:     parser.NewMarker("path:to:c128", parser.MarkerKindComplex, reflect.ValueOf(_complexValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.C128(0 + 0i)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				c := got.Interface().(sdktesting.C128)
				return c == _complexValue
			},
		},
		{
			name:    "complex marker to ptrc64 type",
			mrk:     parser.NewMarker("path:to:ptrc64", parser.MarkerKindComplex, reflect.ValueOf(_complexValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrC64(new(complex64))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				c := got.Interface().(sdktesting.PtrC64)
				return *c == _complexValue
			},
		},
		{
			name:    "complex marker to ptrc128 type",
			mrk:     parser.NewMarker("path:to:ptrc128", parser.MarkerKindComplex, reflect.ValueOf(_complexValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrC128(new(complex128))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				c := got.Interface().(sdktesting.PtrC128)
				return *c == _complexValue
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
