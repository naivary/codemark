package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
)

type str string
type ptrstr *string

func strDefs(t *testing.T) Registry {
	reg := NewInMemoryRegistry()
	defs := []*Definition{
		MakeDef("path:to:str", TargetField, reflect.TypeOf(str(""))),
		MakeDef("path:to:ptrstr", TargetField, reflect.TypeOf(ptrstr(new(string)))),
		MakeDef("path:to:rune", TargetField, reflect.TypeOf(r(0))),
		MakeDef("path:to:ptrrune", TargetField, reflect.TypeOf(ptrrune(new(rune)))),
		MakeDef("path:to:byte", TargetField, reflect.TypeOf(b(0))),
		MakeDef("path:to:ptrbyte", TargetField, reflect.TypeOf(ptrbyte(new(byte)))),
	}
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
			t.Errorf("err occured: %s\n", err)
		}

	}
	return reg
}

func TestStringConverter(t *testing.T) {
	tests := []struct {
		name          string
		mrk           parser.Marker
		t             Target
		out           reflect.Type
		expectedValue any
		isValid       bool
		isValidValue  func(got reflect.Value) bool
	}{
		{
			name:          "string marker to string type",
			mrk:           parser.NewMarker("path:to:str", parser.MarkerKindString, reflect.ValueOf(string("codemark"))),
			t:             TargetField,
			out:           reflect.TypeOf(str("")),
			isValid:       true,
			expectedValue: "codemark",
			isValidValue: func(got reflect.Value) bool {
				str := got.Interface().(str)
				return str == "codemark"
			},
		},
		{
			name:          "string marker to ptr string type",
			mrk:           parser.NewMarker("path:to:ptrstr", parser.MarkerKindString, reflect.ValueOf(string("codemark"))),
			t:             TargetField,
			out:           reflect.TypeOf(ptrstr(new(string))),
			expectedValue: "codemark",
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				str := got.Interface().(ptrstr)
				return *str == "codemark"
			},
		},
		{
			name:          "string marker to rune type",
			mrk:           parser.NewMarker("path:to:rune", parser.MarkerKindString, reflect.ValueOf(string("c"))),
			t:             TargetField,
			out:           reflect.TypeOf(r(0)),
			expectedValue: "c",
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				r := got.Interface().(r)
				return r == 'c'
			},
		},
		{
			name:          "string marker to ptr rune type",
			mrk:           parser.NewMarker("path:to:ptrrune", parser.MarkerKindString, reflect.ValueOf(string("c"))),
			t:             TargetField,
			out:           reflect.TypeOf(ptrrune(new(rune))),
			expectedValue: "c",
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				r := got.Interface().(ptrrune)
				return *r == 'c'
			},
		},
		{
			name:          "string marker to byte type",
			mrk:           parser.NewMarker("path:to:byte", parser.MarkerKindString, reflect.ValueOf(string("c"))),
			t:             TargetField,
			out:           reflect.TypeOf(b(0)),
			expectedValue: "c",
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				b := got.Interface().(b)
				return b == 'c'
			},
		},
		{
			name:          "string marker to ptr byte type",
			mrk:           parser.NewMarker("path:to:ptrbyte", parser.MarkerKindString, reflect.ValueOf(string("c"))),
			t:             TargetField,
			out:           reflect.TypeOf(ptrbyte(new(byte))),
			expectedValue: "c",
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				b := got.Interface().(ptrbyte)
				return *b == 'c'
			},
		},
	}
	reg := strDefs(t)
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
