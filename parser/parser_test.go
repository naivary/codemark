package parser

import (
	"reflect"
	"testing"
)

// TODO: better testing of parser..
func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		isValid bool
		kind    reflect.Kind
	}{
		{
			name:    "boolean without assignment",
			input:   "+jsonschema:validation:required",
			isValid: true,
			kind:    reflect.Bool,
		},
		{
			name:    "int",
			input:   "+jsonschema:validation:max=2",
			isValid: true,
			kind:    reflect.Int64,
		},
		{
			name:    "hex",
			input:   "+jsonschema:validation:max=0x23f3e",
			isValid: true,
			kind:    reflect.Int64,
		},
		{
			name:    "complex",
			input:   "+jsonschema:validation:max=0x23ef+2i",
			isValid: true,
			kind:    reflect.Complex128,
		},
		{
			name:    "catch error",
			input:   "+jsonschema:validation:max=",
			isValid: false,
		},
		{
			name: "multi markers",
			input: `+jsonschema:validation:max=3


+jsonschema:validation:max=5`,
			isValid: true,
		},
		{
			name:    "array",
			input:   `+jsonschema:validation:items=["lorem", "ipsum", "levy"]`,
			isValid: true,
			kind:    reflect.Slice,
		},
		{
			name:    "array escape and int",
			input:   `+jsonschema:validation:items=["lorem", "ips\"um", 3]`,
			isValid: true,
			kind:    reflect.Slice,
		},
		{
			name:    "array escape and float",
			input:   `+jsonschema:validation:items=["lorem", "ips\"um", 3.3]`,
			isValid: true,
			kind:    reflect.Slice,
		},
		{
			name:    "bool with assignment",
			input:   `+jsonschema:validation:required=false`,
			isValid: true,
			kind:    reflect.Bool,
		},
		{
			name:    "array with bool",
			input:   `+jsonschema:validation:items=["lorem", "ips\"um", true]`,
			isValid: true,
			kind:    reflect.Slice,
		},
		{
			name:    "float",
			input:   `+codemark:parser:float=3.2`,
			isValid: true,
			kind:    reflect.Float64,
		},
		{
			name: "multiline string",
			input: `+codemark:parser:string.multiline=` + "`" + `multi line 
string` + "`",
			isValid: true,
			kind:    reflect.String,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			markers, err := Parse(tc.input)
			if err != nil && tc.isValid {
				t.Fatalf("Expected to be valid but got an error marker: %v", err)
			}
			for _, m := range markers {
				t.Logf("kind: `%v`; value: `%v`\n", m.Kind(), m.Value())
			}
		})
	}
}
