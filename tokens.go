package main

import "fmt"

type TokenKind int

const (
	// reached the edn
	TokenKindEOF = iota + 1
	// e.g. +path:to:marker=[item,second,third]
	TokenKindArray
	// e.g. +path:to:marker=string_value
	TokenKindString
	// e.g. +path:to:marker=3
	TokenKindNumber
	// e.g. +path:to:marker. If the token is present the value is always true.
	TokenKindBool
	// e.g. +path:to:marker={key:value, key:value}
	TokenKindMap
	// e.g. path:to:marker
	TokenKindIdent
	// e.g. `=`
	TokenKindAssignment
	// e.g. `{`
	TokenKindOpenCurly
	// e.g. `}`
	TokenKindCloseCurly
	// e.g. `+` before the identifier
	TokenKindPlus
	// e.g. `[`
	TokenKindOpenSquareBracket
	// e.g. `[`
	TokenKindCloseSquareBracket
	// represent the value of an array
	TokenKindArrayValue

    TokenKindDoubleQuotationMark

    TokenKindCommaSeparator

	TokenKindError
)

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
