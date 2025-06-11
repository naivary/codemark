package lexer

import (
	"fmt"
)

type TokenKind int

const (
	// reached the edn
	TokenKindEOF TokenKind = iota + 1
	// e.g. +path:to:marker=somekindofvalue
	TokenKindString
	// e.g. +path:to:marker. If the token is present the value is always true.
	TokenKindBool
	// e.g. path:to:marker
	TokenKindIdent
	// e.g. `=`
	TokenKindAssignment
	// e.g. `+` before the identifier
	TokenKindPlus
	// e.g. `[`
	TokenKindOpenSquareBracket
	// e.g. `]`
	TokenKindCloseSquareBracket
	// e.g. 1237123, 0x283f etc.
	TokenKindInt
	// e.g 1.2
	TokenKindFloat
	// e.g. 3 + 2i
	TokenKindComplex

	TokenKindError
)

var tokenNames = map[TokenKind]string{
	TokenKindEOF:                "TokenKindEOF",
	TokenKindString:             "TokenKindString",
	TokenKindBool:               "TokenKindBool",
	TokenKindIdent:              "TokenKindIdent",
	TokenKindAssignment:         "TokenKindAssignment",
	TokenKindPlus:               "TokenKindPlus",
	TokenKindOpenSquareBracket:  "TokenKindOpenSquareBracket",
	TokenKindCloseSquareBracket: "TokenKindCloseSquareBracket",
	TokenKindInt:                "TokenKindInt",
	TokenKindFloat:              "TokenKindFloat",
	TokenKindComplex:            "TokenKindComplex",
	TokenKindError:              "TokenKindError",
}

func (t TokenKind) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return fmt.Sprintf("TokenKind<%d>", t)
}

func NewToken(kind TokenKind, v string) Token {
	return Token{
		Kind:  kind,
		Value: v,
	}
}

type Token struct {
	Kind  TokenKind
	Value string
}

func (t Token) String() string {
	switch t.Kind {
	case TokenKindEOF:
		return "EOF"
	case TokenKindError:
		return t.Value
	}
	return fmt.Sprintf("%q", t.Value)
}
