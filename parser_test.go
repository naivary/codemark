package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "boolean",
			input: "+jsonschema:validation:required",
		},
		{
			name:  "int",
			input: "+jsonschema:validation:max=2",
		},
		{
			name:  "hex",
			input: "+jsonschema:validation:max=0x23f3e",
		},
		{
			name:  "complex",
			input: "+jsonschema:validation:max=0x23ef+2i",
		},
		{
			name:  "catch error",
			input: "+jsonschema:validation:max=",
		},
		{
			name: "multi markers",
			input: `+jsonschema:validation:max=3


+jsonschema:validation:max=5`,
		},
		{
			name:  "array",
            input: `+jsonschema:validation:items=["lorem", "ipsum", "levy"]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := parse(tc.input)
			if err := p.run(); err != nil {
				t.Fatal(err)
			}
			for m := range p.markers {
				t.Log(m.Kind())
				t.Log(m.Value())
			}
		})
	}
}
