package lexer

import (
	"strings"
	"unicode"
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

func isSingleLineString(r rune) bool {
	return r == _dquot 
}

func isMultiLineString(r rune) bool {
	return r == _tick 
}

func isNewline(r rune) bool {
	return r == _newline || r == _return
}

func hasPlusPrefix(input string, pos int) bool {
	return strings.HasPrefix(input[pos:], string(_plus))
}

// isLower checks if `r` is a lower letter
func isLower(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsLower(r)
}
