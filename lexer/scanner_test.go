package lexer

import (
	"testing"
)

func TestScanString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		isValid bool
	}{
		{
			name:    "normal input",
			input:   `"just a string"`,
			isValid: true,
		},
		{
            name:    "escaped input",
			input:   `"just a string\"more string"`,
			isValid: true,
		},

	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := Lex(tc.input)
			// skip the first double qoute
			l.next()
			l.ignore()

			err := scanString(l)
			if err != nil && tc.isValid {
				t.Fatalf("expected to be a valid input but got an error: `%s`", err.Error())
			}
		})
	}
}
