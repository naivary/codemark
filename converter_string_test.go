package codemark

import (
	"reflect"
	"slices"
	"testing"

	"github.com/naivary/codemark/parser"
)

type str string
type ptrstr *string
type r rune
type ptrrune *rune
type b byte
type ptrbyte *byte
type bytes []byte
type ptrbytes []*byte
type runes []rune
type ptrrunes []*rune

func strDefs(t *testing.T) Registry {
	reg := NewInMemoryRegistry()
	defs := []*Definition{
		MakeDef("path:to:str", TargetField, reflect.TypeOf(str(""))),
		MakeDef("path:to:ptrstr", TargetField, reflect.TypeOf(ptrstr(new(string)))),
		MakeDef("path:to:rune", TargetField, reflect.TypeOf(r(0))),
		MakeDef("path:to:ptrrune", TargetField, reflect.TypeOf(ptrrune(new(rune)))),
		MakeDef("path:to:byte", TargetField, reflect.TypeOf(b(0))),
		MakeDef("path:to:ptrbyte", TargetField, reflect.TypeOf(ptrbyte(new(byte)))),
		MakeDef("path:to:bytes", TargetField, reflect.TypeOf(bytes([]byte{}))),
		MakeDef("path:to:ptrbytes", TargetField, reflect.TypeOf(ptrbytes([]*byte{}))),
		MakeDef("path:to:runes", TargetField, reflect.TypeOf(runes([]rune{}))),
		MakeDef("path:to:ptrrunes", TargetField, reflect.TypeOf(ptrrunes([]*rune{}))),
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
		{
			name:          "string marker to bytes types single value",
			mrk:           parser.NewMarker("path:to:bytes", parser.MarkerKindString, reflect.ValueOf(string("c"))),
			t:             TargetField,
			out:           reflect.TypeOf(bytes([]byte{})),
			expectedValue: []byte("c"),
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				b := got.Interface().(bytes)
				want := []byte("c")
				return slices.Equal(b, want)
			},
		},
		{
			name:          "string marker to bytes types char chain",
			mrk:           parser.NewMarker("path:to:bytes", parser.MarkerKindString, reflect.ValueOf(string("codemark"))),
			t:             TargetField,
			out:           reflect.TypeOf(bytes([]byte{})),
			expectedValue: []byte("codemark"),
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				b := got.Interface().(bytes)
				want := []byte("codemark")
				return slices.Equal(b, want)
			},
		},
		{
			name:          "string marker to ptr bytes types single value",
			mrk:           parser.NewMarker("path:to:ptrbytes", parser.MarkerKindString, reflect.ValueOf(string("c"))),
			t:             TargetField,
			out:           reflect.TypeOf(ptrbytes([]*byte{})),
			expectedValue: []byte("c"),
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				bytes := got.Interface().(ptrbytes)
				want := []byte("c")
				for i, b := range bytes {
					if *b != want[i] {
						return false
					}
				}
				return true
			},
		},
		{
			name:          "string marker to ptr bytes types char chain",
			mrk:           parser.NewMarker("path:to:ptrbytes", parser.MarkerKindString, reflect.ValueOf(string("codemark"))),
			t:             TargetField,
			out:           reflect.TypeOf(ptrbytes([]*byte{})),
			expectedValue: []byte("codemark"),
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				bytes := got.Interface().(ptrbytes)
				want := []byte("codemark")
				for i, b := range bytes {
					if *b != want[i] {
						return false
					}
				}
				return true
			},
		},
		{
			name:          "string marker to runes type",
			mrk:           parser.NewMarker("path:to:runes", parser.MarkerKindString, reflect.ValueOf(string("codemark"))),
			t:             TargetField,
			out:           reflect.TypeOf(runes([]rune{})),
			expectedValue: []rune("codemark"),
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				runes := got.Interface().(runes)
				want := []rune("codemark")
				for i, r := range runes {
					if r != want[i] {
						return false
					}
				}
				return true
			},
		},
		{
			name:          "string marker to ptr runes type",
			mrk:           parser.NewMarker("path:to:ptrrunes", parser.MarkerKindString, reflect.ValueOf(string("codemark"))),
			t:             TargetField,
			out:           reflect.TypeOf(ptrrunes([]*rune{})),
			expectedValue: []rune("codemark"),
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				runes := got.Interface().(ptrrunes)
				want := []rune("codemark")
				for i, r := range runes {
					if *r != want[i] {
						return false
					}
				}
				return true
			},
		},
	}
	reg := strDefs(t)
	mngr, err := NewConvMngr(reg, &stringConverter{})
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
