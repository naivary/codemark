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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := parse(tc.input)
			p.run()
			for m := range p.markers {
				t.Log(m)
			}
		})
	}
}
