package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
)

const _boolValue = true

type boolean bool
type ptrboolean *bool

func boolDefs(t *testing.T) Registry {
	reg := NewInMemoryRegistry()
	defs := []*Definition{
		// bool
		MakeDef("path:to:bool", TargetField, reflect.TypeOf(boolean(false))),
		// ptr bool
		MakeDef("path:to:ptrbool", TargetField, reflect.TypeOf(ptrboolean(new(bool)))),
	}
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
			t.Errorf("err occured: %s\n", err)
		}

	}
	return reg
}

func TestBoolConverter(t *testing.T) {
	tests := []struct {
		name         string
		mrk          parser.Marker
		t            Target
		out          reflect.Type
		isValid      bool
		isValidValue func(got reflect.Value) bool
	}{
		{
			name:    "bool marker to bool type",
			mrk:     parser.NewMarker("path:to:bool", parser.MarkerKindBool, reflect.ValueOf(_boolValue)),
			t:       TargetField,
			out:     reflect.TypeOf(boolean(false)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				b := got.Interface().(boolean)
				return b == _boolValue
			},
		},
		{
			name:    "bool marker to ptr bool type",
			mrk:     parser.NewMarker("path:to:ptrbool", parser.MarkerKindBool, reflect.ValueOf(_boolValue)),
			t:       TargetField,
			out:     reflect.TypeOf(ptrboolean(new(bool))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				b := got.Interface().(ptrboolean)
				return *b == _boolValue
			},
		},
	}
	reg := boolDefs(t)
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
