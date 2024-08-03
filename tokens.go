package main

import "fmt"

type TokenKind int

const (
	// reached the edn
	TokenKindEOF TokenKind = iota + 1
	// e.g. +path:to:marker=[item,second,third]
	TokenKindArray
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
	// e.g. `[`
	TokenKindCloseSquareBracket
	// e.g. 1237123, 0x283f etc.
	TokenKindInt
	// e.g 1.2
	TokenKindFloat
	// e.g. 3 + 2i
	TokenKindComplex

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
