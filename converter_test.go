package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
)

func ptr[T any](t T) *T {
	return &t
}

type str string
type strPtr *string

type i int
type i8 int8
type i16 int16
type i32 int32
type i64 int64

type b byte
type r rune

type byteSlice []byte
type runeSlice []rune

type boolean bool
type boolPtr bool

func TestConverter_Convert(t *testing.T) {
	tests := []struct {
		name    string
		markers []parser.Marker
		defs    []*Definition
	}{
		{
			name: "string marker",
			markers: []parser.Marker{
				parser.NewMarker("jsonschema:validation:string", reflect.String, reflect.ValueOf("string")),
				parser.NewMarker("jsonschema:validation:stringptr", reflect.String, reflect.ValueOf(ptr("stringptr"))),
				parser.NewMarker("jsonschema:validation:byte", reflect.String, reflect.ValueOf("b")),
				parser.NewMarker("jsonschema:validation:rune", reflect.String, reflect.ValueOf("r")),
				parser.NewMarker("jsonschema:validation:runeslice", reflect.String, reflect.ValueOf("rrrrrr")),
				parser.NewMarker("jsonschema:validation:byteslice", reflect.String, reflect.ValueOf("bbbbbb")),
			},
			defs: []*Definition{
				MakeDef("jsonschema:validation:string", TargetConst, reflect.TypeOf(str(""))),
				MakeDef("jsonschema:validation:stringptr", TargetConst, reflect.TypeOf(strPtr(ptr("")))),
				MakeDef("jsonschema:validation:byte", TargetConst, reflect.TypeOf(b(0))),
				MakeDef("jsonschema:validation:rune", TargetConst, reflect.TypeOf(r(0))),
				MakeDef("jsonschema:validation:byteslice", TargetConst, reflect.TypeOf(byteSlice([]byte{}))),
				MakeDef("jsonschema:validation:runeslice", TargetConst, reflect.TypeOf(runeSlice([]rune{}))),
			},
		},
		{
			name: "boolean marker",
			markers: []parser.Marker{
				parser.NewMarker("jsonschema:validation:required", reflect.Bool, reflect.ValueOf(true)),
				parser.NewMarker("jsonschema:validation:requiredptr", reflect.Bool, reflect.ValueOf(true)),
			},
			defs: []*Definition{
				MakeDef("jsonschema:validation:required", TargetConst, reflect.TypeOf(boolean(false))),
				MakeDef("jsonschema:validation:requiredptr", TargetConst, reflect.TypeOf(boolPtr(false))),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reg := NewRegistry()
			for _, def := range tc.defs {
				if err := reg.Define(def); err != nil {
					t.Error(err)
				}
			}
			conv, err := NewConverter(reg)
			if err != nil {
				t.Error(err)
			}
			for _, marker := range tc.markers {
				v, err := conv.Convert(marker)
				if err != nil {
					t.Error(err)
				}
				fmt.Println(reflect.TypeOf(v))
			}

		})
	}

}
