package main

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/marker"
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

type iPtr *int
type i8Ptr *int8
type i16Ptr *int16
type i32Ptr *int32
type i64Ptr *int64

type b byte
type bytePtr *byte
type r rune
type runePtr *rune

type byteSlice []byte
type byteSlicePtr *[]byte

type runeSlice []rune
type runeSlicePtr *[]rune

type boolean bool
type boolPtr *bool

func TestConverter_Convert(t *testing.T) {
	tests := []struct {
		name    string
		markers []marker.Marker
		defs    []*Definition
	}{
		{
			name: "string marker",
			markers: []marker.Marker{
				marker.NewDefault("jsonschema:validation:string", reflect.String, reflect.ValueOf("string")),
				marker.NewDefault("jsonschema:validation:byte", reflect.String, reflect.ValueOf("b")),
				marker.NewDefault("jsonschema:validation:rune", reflect.String, reflect.ValueOf("r")),
				marker.NewDefault("jsonschema:validation:runeslice", reflect.String, reflect.ValueOf("rrrrrr")),
				marker.NewDefault("jsonschema:validation:byteslice", reflect.String, reflect.ValueOf("bbbbbb")),
			},
			defs: []*Definition{
				MakeDef("jsonschema:validation:string", TargetConst, reflect.TypeOf(str(""))),
				MakeDef("jsonschema:validation:byte", TargetConst, reflect.TypeOf(b(0))),
				MakeDef("jsonschema:validation:rune", TargetConst, reflect.TypeOf(r(0))),
				MakeDef("jsonschema:validation:byteslice", TargetConst, reflect.TypeOf(byteSlice([]byte{}))),
				MakeDef("jsonschema:validation:runeslice", TargetConst, reflect.TypeOf(runeSlice([]rune{}))),
			},
		},
		{
			name: "boolean marker",
			markers: []marker.Marker{
				marker.NewDefault("jsonschema:validation:required", reflect.Bool, reflect.ValueOf(true)),
				marker.NewDefault("jsonschema:validation:requiredptr", reflect.Bool, reflect.ValueOf(true)),
			},
			defs: []*Definition{
				MakeDef("jsonschema:validation:required", TargetConst, reflect.TypeOf(boolean(false))),
				MakeDef("jsonschema:validation:requiredptr", TargetConst, reflect.TypeOf(boolPtr(ptr(false)))),
			},
		},
		{
			name: "string marker pointers",
			markers: []marker.Marker{
				marker.NewDefault("jsonschema:validation:stringptr", reflect.String, reflect.ValueOf("stringptr")),
				marker.NewDefault("jsonschema:validation:byteptr", reflect.String, reflect.ValueOf("b")),
				marker.NewDefault("jsonschema:validation:runeptr", reflect.String, reflect.ValueOf("r")),
			},
			defs: []*Definition{
				MakeDef("jsonschema:validation:stringptr", TargetConst, reflect.TypeOf(strPtr(ptr("")))),
				MakeDef("jsonschema:validation:byteptr", TargetConst, reflect.TypeOf(bytePtr(ptr[byte](0)))),
				MakeDef("jsonschema:validation:runeptr", TargetConst, reflect.TypeOf(runePtr(ptr[rune](0)))),
			},
		},
		{
			name: "integer marker",
			markers: []marker.Marker{
				marker.NewDefault("jsonschema:validation:i", reflect.Int64, reflect.ValueOf(30)),
				marker.NewDefault("jsonschema:validation:i8", reflect.Int64, reflect.ValueOf(2)),
			},
			defs: []*Definition{
				MakeDef("jsonschema:validation:i", TargetConst, reflect.TypeOf(i(0))),
				MakeDef("jsonschema:validation:i8", TargetConst, reflect.TypeOf(i8(0))),
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
				v, err := conv.Convert(marker, TargetConst)
				if err != nil {
					t.Error(err)
				}
				t.Log(v)
			}

		})
	}

}
