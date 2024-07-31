package main

import (
	"testing"
)

func TestLexer_Lex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		isValid bool
		want    bool
	}{
		{
			name:    "random input before",
			input:   "dasdasd asd asdsadasd as +jsonschema:validation:maximum=3",
			isValid: true,
			want:    true,
		},
		{
			name:    "number input",
			input:   "+jsonschema:validation:maximum=3",
			isValid: true,
			want:    true,
		},
		{
			name:    "boolean input",
			input:   "+jsonschema:validation:required",
			isValid: true,
			want:    true,
		},
		{
			name:    "string input",
			input:   "+jsonschema:validation:format=email",
			isValid: true,
			want:    true,
		},
        {
            name: "string with plus",
            input: "asdhajsdhjds+dhsajdhasjhdashdjad +jsonschema:validation:maximum=3",
            isValid: true,
            want: true,
        },
        {
            name: "multi line string",
            input: `dashjdhsajdh jasdhjasdh jasdhjashdasjdhasjdhj

jsonschema:validation:maximum=3`,
            isValid: true,
            want: true,
        },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := lex(tc.input)
			l.run()
			for token := range l.tokens {
				t.Log(token)
			}
		})
	}
}
