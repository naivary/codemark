package lexer

import (
	"fmt"

	"github.com/naivary/codemark/internal/lexer/token"
)

func NewToken(kind token.Kind, v string) Token {
	return Token{
		Kind:  kind,
		Value: v,
	}
}

type Token struct {
	Kind  token.Kind
	Value string
}

func (t Token) String() string {
	switch t.Kind {
	case token.EOF:
		return "EOF"
	case token.ERROR:
		return t.Value
	}
	return fmt.Sprintf("%q", t.Value)
}
