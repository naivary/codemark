package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
)

type stringList []string
type ptrStringList []*string

type intList []int
type ptrIntList []*int

func listDefs(t *testing.T) Registry {
	reg := NewInMemoryRegistry()
	defs := []*Definition{
		// string list
		MakeDef("path:to:stringlist", TargetField, reflect.TypeOf(stringList([]string{}))),
		MakeDef("path:to:intlist", TargetField, reflect.TypeOf(intList([]int{}))),

		// ptr
		MakeDef("path:to:ptrstringlist", TargetField, reflect.TypeOf(ptrStringList([]*string{}))),
		MakeDef("path:to:ptrintlist", TargetField, reflect.TypeOf(ptrIntList([]*int{}))),
	}
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
			t.Errorf("err occured: %s\n", err)
		}

	}
	return reg
}

func TestListConverter(t *testing.T) {
	tests := []struct {
		name         string
		mrk          parser.Marker
		t            Target
		out          reflect.Type
		value        []any
		isValid      bool
		isValidValue func(got reflect.Value, expected []any) bool
	}{
		// string list
		{
			name:    "list marker to string list",
			mrk:     parser.NewMarker("path:to:stringlist", parser.MarkerKindList, reflect.ValueOf([]any{"path", "to", "marker"})),
			t:       TargetField,
			out:     reflect.TypeOf(stringList([]string{})),
			value:   []any{"path", "to", "marker"},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(stringList)
				for i, el := range list {
					expectedElem := expected[i].(string)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},
		{
			name:    "list marker to int list",
			mrk:     parser.NewMarker("path:to:intlist", parser.MarkerKindList, reflect.ValueOf([]any{2, 3, 4})),
			t:       TargetField,
			out:     reflect.TypeOf(intList([]int{})),
			value:   []any{2, 3, 4},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(intList)
				for i, el := range list {
					expectedElem := expected[i].(int)
					if expectedElem != el {
						return false
					}
				}
				return true
			},
		},

		// pointer
		{
			name:    "list marker to ptr string list",
			mrk:     parser.NewMarker("path:to:ptrstringlist", parser.MarkerKindList, reflect.ValueOf([]any{"path", "to", "marker"})),
			t:       TargetField,
			out:     reflect.TypeOf(ptrStringList([]*string{})),
			value:   []any{"path", "to", "marker"},
			isValid: true,
			isValidValue: func(got reflect.Value, expected []any) bool {
				list := got.Interface().(ptrStringList)
				for i, el := range list {
					expectedElem := expected[i].(string)
					if expectedElem != *el {
						return false
					}
				}
				return true
			},
		},
	}
	reg := listDefs(t)
	mngr, err := NewConvMngr(reg, &listConverter{})
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
			if !tc.isValidValue(rvalue, tc.value) {
				t.Fatalf("value is not correct. got: %v", rvalue)
			}
			t.Logf("succesfully converted. got: %v; expected: %v\n", rtype, tc.out)
		})
	}
}
