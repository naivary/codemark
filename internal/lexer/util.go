package lexer

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/naivary/codemark/internal/lexer/token"
	"github.com/naivary/codemark/internal/utils"
)

func ignoreSpace(l *Lexer) {
	l.acceptFunc(isSpace)
	l.ignore()
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r) || r == '-' || r == '+'
}

func isSpace(r rune) bool {
	return r == _whitespace || r == _tab
}

func isBool(r rune) bool {
	return r == 't' || r == 'f'
}

func isNewline(r rune) bool {
	return r == _newline || r == _return
}

func hasPlusPrefix(input string, pos int) bool {
	return strings.HasPrefix(input[pos:], string(_plus))
}

func isLower(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsLower(r)
}

func kindOfNumber(number string) (token.Kind, error) {
	kind := token.INT
	_, errInt := strconv.ParseInt(number, 0, 64)
	_, errFloat := strconv.ParseFloat(number, 64)
	c, errComplex := strconv.ParseComplex(utils.ComplexOrder(number), 128)

	if errInt != nil && errFloat != nil && errComplex != nil {
		return token.ERROR, fmt.Errorf("cannot lex number: %s", number)
	}
	if errInt != nil && errFloat == nil {
		kind = token.FLOAT
	}
	if errComplex == nil && imag(c) != 0 {
		kind = token.COMPLEX
	}
	return kind, nil
}
