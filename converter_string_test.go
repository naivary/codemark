package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
)

type str string
type ptrstr *string
type r rune

func testDefs(t *testing.T) Registry {
	reg := NewInMemoryRegistry()
	defs := []*Definition{
		MakeDef("path:to:str", TargetField, reflect.TypeOf(str(""))),
		MakeDef("path:to:ptrstr", TargetField, reflect.TypeOf(ptrstr(new(string)))),
		MakeDef("path:to:rune", TargetField, reflect.TypeOf(r(0))),
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
		name string
		mrk  parser.Marker
		t    Target
		out  reflect.Type
	}{
		{
			name: "string marker to string type",
			mrk:  parser.NewMarker("path:to:str", parser.MarkerKindString, reflect.ValueOf(string("codemark"))),
			t:    TargetField,
			out:  reflect.TypeOf(str("")),
		},
		{
			name: "string marker to ptr string type",
			mrk:  parser.NewMarker("path:to:ptrstr", parser.MarkerKindString, reflect.ValueOf(string("codemark"))),
			t:    TargetField,
			out:  reflect.TypeOf(ptrstr(new(string))),
		},
		{
			name: "string marker to rune type",
			mrk:  parser.NewMarker("path:to:rune", parser.MarkerKindString, reflect.ValueOf(string("c"))),
			t:    TargetField,
			out:  reflect.TypeOf(r(0)),
		},
	}
	reg := testDefs(t)
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
			t.Logf("succesfully converted. got: %v; expected: %v\n", rtype, tc.out)
		})
	}
}
