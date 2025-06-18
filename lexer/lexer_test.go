package lexer

import (
	"testing"

	"github.com/naivary/codemark/lexer/token"
)

func TestLexer_Lex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		isValid bool
	}{
		{
			name:    "random input before",
			input:   `dasdasd asd asdsadasd as +jsonschema:validation:maximum=3`,
			isValid: true,
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
			input:   `+jsonschema:validation:format="email"`,
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
+jsonschema:validation:format="email"`,
			isValid: true,
		},
		{
			name: "multi marker without doc reversed",
			input: `
+jsonschema:validation:maximum=3
+jsonschema:validatiom:format="email"
this is a normal docs string`,
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
+jsonschema:validatiom:format="email"`,
			isValid: false,
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
			input:   `+jsonschema:validation:items=["somet3s"]`,
			isValid: true,
		},
		{
			name:    "space between end bracket",
			input:   `+jsonschema:validation:items=[somet3s   ]`,
			isValid: false,
		},
		{
			name:    "escaped",
			input:   `+jsonschema:validation:items=["something\""]`,
			isValid: true,
		},
		{
			name:    "escaped with letter in between",
			input:   `+jsonschema:validation:items=["something\"s"]`,
			isValid: true,
		},
		{
			name:    "invalid ending because its followed by space",
			input:   `+jsonschema:validation:items=["something\"s" ]`,
			isValid: false,
		},
		{
			name:    "complex number",
			input:   `+jsonschema:validation:max=2i+3`,
			isValid: false,
		},
		{
			name:    "complex number valid",
			input:   `+jsonschema:validation:max=3+2i`,
			isValid: true,
		},
		{
			name:    "unfinished assignment",
			input:   `+jsonschema:validation:max=`,
			isValid: false,
		},
		{
			name:    "float",
			input:   `+jsonschema:validation:max=3.5`,
			isValid: true,
		},
		{
			name:    "bool with assignment",
			input:   `+jsonschema:validation:max=true`,
			isValid: true,
		},
		{
			name:    "bool with assignment false",
			input:   `+jsonschema:validation:max=false`,
			isValid: true,
		},
		{
			name:    "bool with assignment false with spaces and number after",
			input:   `+jsonschema:validation:max=false      3`,
			isValid: false,
		},
		{
			name:    "negative integer",
			input:   `+jsonschema:validation:max=-3`,
			isValid: true,
		},
		{
			name:    "array with boolean",
			input:   `+jsonschema:validation:max=["name", true]`,
			isValid: true,
		},
		{
			name:    "array start with bool",
			input:   `+jsonschema:validation:max=[true, true]`,
			isValid: true,
		},
		{
			name:    "array with all possible types",
			input:   `+jsonschema:validation:max=[true, false, "some-string", 2, 0x24, 3.21, 3+2i, -2]`,
			isValid: true,
		},
		{
			name:    "invalid string value",
			input:   `+jsonschema:validation:name="name\"`,
			isValid: false,
		},
		{
			name:    "invalid identifier",
			input:   `+jsonschema:validation="name"`,
			isValid: false,
		},
		{
			name: "boolean without assignment wiht doc",
			input: `+jsonschema:validation:required
this is the doc`,
			isValid: true,
		},
		{
			name:    "identifier is not allowed to start with a number",
			input:   `+3idn:validation:required`,
			isValid: false,
		},
		{
			name:    "identifier is allowed to include numbers",
			input:   `+idn3:validation:required`,
			isValid: true,
		},
		{
			name:    "identifier is allowed to include underscore",
			input:   `+idn_v3:validation:required`,
			isValid: true,
		},
		{
			name:    "identifier is allowed to include dots",
			input:   `+idn_v3.1:validation:required`,
			isValid: true,
		},
		{
			name:    "identifier is not allowed to end with underscore (end segment)",
			input:   `+idn_v3.1:validation:required_`,
			isValid: false,
		},
		{
			name:    "identifier is not allowed to end with underscore (middle segment)",
			input:   `+idn_v3.1:validation_:required`,
			isValid: false,
		},
		{
			name:    "identifier is not allowed to end with underscore (first segment)",
			input:   `+idn_v3.1_:validation:required`,
			isValid: false,
		},
		{
			name:    "identifier is not allowed to end with dot (end segment)",
			input:   `+idn_v3.1:validation:required.`,
			isValid: false,
		},
		{
			name:    "identifier is not allowed to end with dot (middle segment)",
			input:   `+idn_v3.1:validation.:required`,
			isValid: false,
		},
		{
			name:    "identifier is not allowed to end with dot (first segment)",
			input:   `+idn_v3.1.:validation:required`,
			isValid: false,
		},
		{
			name: "newline followed by EOF",
			input: `+idn_v3.1:validation:required
`,
			isValid: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := Lex(tc.input)
			for tk := range l.tokens {
				t.Log(tk)
				if tk.Kind == token.ERROR && tc.isValid {
					t.Fatalf("expected to lex correctly: `%s`. Error is: %s", tc.input, tk.Value)
				}

				if tk.Kind == token.EOF && !tc.isValid {
					t.Fatalf("expected an error and didn't get one for `%s`", tc.input)
				}
			}
		})
	}
}
