package main

import (
	"testing"
)

func TestLexer_Lex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		isValid bool
	}{
		{
			name: "random input before",
			input: `dasdasd asd asdsadasd as 
            +jsonschema:validation:maximum=3`,
			isValid: false,
		},
		{
			name:    "number input",
			input:   "+jsonschema:validation:maximum=3",
			isValid: true,
		},
		{
			name:    "boolean input",
			input:   "+jsonschema:validation:required",
			isValid: true,
		},
		{
			name:    "string input",
			input:   "+jsonschema:validation:format=email",
			isValid: true,
		},
		{
			name: "string with plus",
			input: `asdhajsdhjds+dhsajdhasjhdashdjad 

+jsonschema:validation:maximum=3`,
			isValid: true,
		},
		{
			name: "multi line string",
			input: `dashjdhsajdh jasdhjasdh jasdhjashdasjdhasjdhj

+jsonschema:validation:maximum=3`,
			isValid: true,
		},
		{
			name: "multi marker without doc",
			input: `+jsonschema:validation:maximum=3
+jsonschema:validatiom:format=email`,
			isValid: true,
		},
		{
			name: "multi marker without doc reversed",
			input: `+jsonschema:validation:maximum=3
+jsonschema:validatiom:format=email
dasdasd asdasd as dasdasd ads adsd`,
			isValid: true,
		},
		{
			name: "string one newline",
			input: `asdhajsdhjdsdhsajdhasjhdashdjad 
+jsonschema:validation:maximum=3`,
			isValid: true,
		},
		{
			name: "string multiple new lines",
			input: `asdhajsdhjdsdhsajdhasjhdashdjad 





+jsonschema:validation:maximum=3`,
			isValid: true,
		},
		{
			name:    "no identifier",
			input:   `+=3`,
			isValid: false,
		},
		{
			name:    "spaces before line",
			input:   `       +jsonschema:validation:maximum=3`,
			isValid: true,
		},
		{
			name: "multi marker without doc",
			input: `+jsonschema:validation:maximum=3 dashdasjdhasjdhasd
+jsonschema:validatiom:format=email`,
			isValid: true,
		},
		{
			name:    "array with single integer",
			input:   `+jsonschema:validation:items=[21323]`,
			isValid: true,
		},
		{
			name:    "array multi integer",
			input:   `+jsonschema:validation:items=[12345, 1234, 12]`,
			isValid: true,
		},
		{
			name:    "array string",
			input:   `+jsonschema:validation:items=[somet3s]`,
			isValid: true,
		},
		{
			name:    "space between end bracket",
			input:   `+jsonschema:validation:items=[somet3s   ]`,
			isValid: false,
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
