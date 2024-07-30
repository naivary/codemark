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

    TokenKindCommaSeparator

	TokenKindError
)

func (t TokenKind) String() string {
	switch t {
    case TokenKindEOF:
        return "EOF"
	case TokenKindArray:
		return "array"
	case TokenKindString:
		return "string"
	case TokenKindNumber:
		return "number"
	case TokenKindBool:
		return "boolean"
	case TokenKindMap:
		return "map"
	case TokenKindIdent:
		return "identifier"
	case TokenKindAssignment:
		return "assignment"
	case TokenKindOpenCurly:
		return "open_curly"
	case TokenKindCloseCurly:
		return "close_curly"
	case TokenKindPlus:
		return "plus"
	case TokenKindOpenSquareBracket:
		return "open_squared_bracket"
	case TokenKindCloseSquareBracket:
		return "close_squared_bracket"
    case TokenKindArrayValue:
        return "array_value"
    case TokenKindCommaSeparator:
        return "comma"
	case TokenKindError:
		return "error"
	default:
		return "INVALID_KIND"
	}
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
