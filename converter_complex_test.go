package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
)

const _complexValue = 9 + 9i

type c64 complex64
type c128 complex128

type ptrc64 *complex64
type ptrc128 *complex128

func complexDefs(t *testing.T) sdk.Registry {
	reg := NewInMemoryRegistry()
	defs := []*sdk.Definition{
		// complex
		MustMakeDef("path:to:c64", sdk.TargetField, reflect.TypeOf(c64(0+0i))),
		MustMakeDef("path:to:c128", sdk.TargetField, reflect.TypeOf(c128(0+0i))),
		// ptr complex
		MustMakeDef("path:to:ptrc64", sdk.TargetField, reflect.TypeOf(ptrc64(new(complex64)))),
		MustMakeDef("path:to:ptrc128", sdk.TargetField, reflect.TypeOf(ptrc128(new(complex128)))),
	}
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
			t.Errorf("err occured: %s\n", err)
		}

	}
	return reg
}

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
			out:     reflect.TypeOf(c64(0 + 0i)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				c := got.Interface().(c64)
				return c == _complexValue
			},
		},
		{
			name:    "complex marker to c128 type",
			mrk:     parser.NewMarker("path:to:c128", parser.MarkerKindComplex, reflect.ValueOf(_complexValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(c128(0 + 0i)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				c := got.Interface().(c128)
				return c == _complexValue
			},
		},
		{
			name:    "complex marker to ptrc64 type",
			mrk:     parser.NewMarker("path:to:ptrc64", parser.MarkerKindComplex, reflect.ValueOf(_complexValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrc64(new(complex64))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				c := got.Interface().(ptrc64)
				return *c == _complexValue
			},
		},
		{
			name:    "complex marker to ptrc128 type",
			mrk:     parser.NewMarker("path:to:ptrc128", parser.MarkerKindComplex, reflect.ValueOf(_complexValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(ptrc128(new(complex128))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				c := got.Interface().(ptrc128)
				return *c == _complexValue
			},
		},
	}
	reg := complexDefs(t)
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
