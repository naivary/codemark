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
		{
			name:    "unclosed double",
			input:   `"just a string\"`,
			isValid: false,
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

func TestScanNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		isValid bool
		kind    TokenKind
	}{
		{
			name:    "real positive number",
			input:   "3",
			isValid: true,
			kind:    TokenKindInt,
		},
		{
			name:    "real positive number with sign",
			input:   "+3",
			isValid: true,
			kind:    TokenKindInt,
		},
		{
			name:    "negative real number",
			input:   "-3",
			isValid: true,
			kind:    TokenKindInt,
		},
		{
			name:    "positive float number",
			input:   "3.02",
			isValid: true,
			kind:    TokenKindFloat,
		},
		{
			name:    "negative float number",
			input:   "-3.02",
			isValid: true,
			kind:    TokenKindFloat,
		},
		{
			name:    "positive complex number",
			input:   "3+2i",
			isValid: true,
			kind:    TokenKindComplex,
		},
		{
			name:    "negative complex number",
			input:   "-3-2i",
			isValid: true,
			kind:    TokenKindComplex,
		},
		{
			name:    "complex number wrong order",
			input:   "2i+3",
			isValid: false,
			kind:    TokenKindError,
		},
		{
			name:    "big positive int number",
			input:   "293123901273",
			isValid: true,
			kind:    TokenKindInt,
		},
		{
			name:    "big positive float number",
			input:   "293123901273.2831739",
			isValid: true,
			kind:    TokenKindFloat,
		},
	}

	for _, tc := range tests {
		l := Lex(tc.input)
		t.Run(tc.name, func(t *testing.T) {
			kind, err := scanNumber(l)
			if err != nil && tc.isValid {
				t.Fatalf("Expected to be valid. Got an error: %s", err.Error())
			}
			if tc.kind != kind {
				t.Fatalf("Kinds are not equal. Expected `%d` got `%d`", tc.kind, kind)
			}
		})
	}
}
