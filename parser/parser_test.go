package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		isValid bool
	}{
		{
			name:    "boolean",
			input:   "+jsonschema:validation:required",
			isValid: true,
		},
		{
			name:    "int",
			input:   "+jsonschema:validation:max=2",
			isValid: true,
		},
		{
			name:    "hex",
			input:   "+jsonschema:validation:max=0x23f3e",
			isValid: true,
		},
		{
			name:    "complex",
			input:   "+jsonschema:validation:max=0x23ef+2i",
			isValid: true,
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
		},
		{
			name:    "array escape and int",
			input:   `+jsonschema:validation:items=["lorem", "ips\"um", 3]`,
			isValid: true,
		},
		{
			name:    "array escape and float",
			input:   `+jsonschema:validation:items=["lorem", "ips\"um", 3.3]`,
			isValid: true,
		},
		{
			name:    "bool with assignment",
			input:   `+jsonschema:validation:required=false`,
			isValid: true,
		},
		{
			name:    "array with bool",
			input:   `+jsonschema:validation:items=["lorem", "ips\"um", true]`,
			isValid: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := parse(tc.input)
			p.run()
			/*if err != nil && !tc.isValid {
				return
			}
			if err != nil && tc.isValid {
				t.Fatalf("expected to parse correctly but got error: `%s`", err.Error())
			}
			if err == nil && !tc.isValid {
				t.Fatalf("expected an error. Got nil")
			}*/
			for m := range p.markers {
				t.Logf("kind: `%v`; value: `%v`; isError: `%v`\n", m.Kind(), m.Value(), m.IsError())
			}
		})
	}
}
