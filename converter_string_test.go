package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
)

type str string
type ptrstr *string

func strDefs(t *testing.T) sdk.Registry {
	reg := NewInMemoryRegistry()
	defs := []*sdk.Definition{
		MustMakeDef("path:to:str", sdk.TargetField, reflect.TypeOf(str(""))),
		MustMakeDef("path:to:ptrstr", sdk.TargetField, reflect.TypeOf(ptrstr(new(string)))),
		MustMakeDef("path:to:rune", sdk.TargetField, reflect.TypeOf(r(0))),
		MustMakeDef("path:to:ptrrune", sdk.TargetField, reflect.TypeOf(ptrrune(new(rune)))),
		MustMakeDef("path:to:byte", sdk.TargetField, reflect.TypeOf(b(0))),
		MustMakeDef("path:to:ptrbyte", sdk.TargetField, reflect.TypeOf(ptrbyte(new(byte)))),
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
		t             sdk.Target
		out           reflect.Type
		expectedValue any
		isValid       bool
		isValidValue  func(got reflect.Value) bool
	}{
		{
			name:          "string marker to string type",
			mrk:           parser.NewMarker("path:to:str", parser.MarkerKindString, reflect.ValueOf(string("codemark"))),
			t:             sdk.TargetField,
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
			t:             sdk.TargetField,
			out:           reflect.TypeOf(ptrstr(new(string))),
			expectedValue: "codemark",
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				str := got.Interface().(ptrstr)
				return *str == "codemark"
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
