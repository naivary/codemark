package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
)

type i int
type str string

var strconv = &stringConverter{}

func TestNewConvMngr(t *testing.T) {
	tests := []struct {
		name  string
		conv  Converter
		types []reflect.Type
	}{
		{
			name:  "adding converter",
			conv:  strconv,
			types: strconv.SupportedTypes(),
		},
	}

	reg := NewInMemoryRegistry()
	def := MakeDef("path:to:marker", TargetField, reflect.TypeOf(bool(false)))
	if err := reg.Define(def); err != nil {
		t.Errorf("err occured: %s\n", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mngr, err := NewConvMngr(reg, tc.conv)
			if err != nil {
				t.Errorf("err occured: %s\n", err)
			}
			for _, rtype := range tc.types {
				typeID, err := TypeID(rtype)
				if err != nil {
					t.Errorf("err occured: %s\n", err)
				}
				conv, err := mngr.GetConverter(rtype)
				if err != nil {
					t.Errorf("err occured: %s\n", err)
				}
				if conv == nil || conv != strconv {
					t.Fatalf("converter for type ID `%s` must exist but is not", typeID)
				}
			}
		})
	}
}

func TestConverterManager_Convert(t *testing.T) {
	tests := []struct {
		name string
		mrk  parser.Marker
		t    Target
		out  reflect.Type
	}{
		{
			name: "string marker to str type",
			mrk:  parser.NewMarker("path:to:str", parser.MarkerKindString, reflect.ValueOf(string("codemark"))),
			t:    TargetField,
			out:  reflect.TypeOf(str("")),
		},
		// test wrong marker kind with right converter e.g. MarkerKindString
		// with intConverter
	}

	reg := NewInMemoryRegistry()
	def := MakeDef("path:to:str", TargetField, reflect.TypeOf(str("")))
	if err := reg.Define(def); err != nil {
		t.Errorf("err occured: %s\n", err)
	}
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
			if reflect.TypeOf(v) != tc.out {
				t.Fatalf("types don't match after conversion. got: %v; expected: %v\n", v, tc.out)
			}
		})
	}
}
