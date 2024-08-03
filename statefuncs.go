package main

import (
	"strings"
	"unicode"
)

type validFunc func(rune) bool

func ignoreSpace(l *lexer) {
	l.acceptFunc(isSpace)
	l.ignore()
}

func ignoreNewline(l *lexer) {
	l.acceptFunc(isNewline)
	l.ignore()
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isNewline(r rune) bool {
	return r == '\n' || r == '\r'
}

func hasPlusPrefix(input string, pos int) bool {
	return strings.HasPrefix(input[pos:], plus)
}

// isAlphaLower checks if `r` is lower and a letter
func isAlphaLower(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsLower(r)
}

// isAlphaNumeric is checking if the rune is a letter, digit or underscore
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

type stateFunc func(*lexer) stateFunc

func lexText(l *lexer) stateFunc {
	// check for the position 0
	if hasPlusPrefix(l.input, l.pos) {
		return lexPlus
	}

	for {
		r := l.next()
		if r == eof {
			break
		}
		if isNewline(r) {
			l.acceptFunc(isSpace)
			l.ignore()
			if hasPlusPrefix(l.input, l.pos) {
				return lexPlus
			}
		}
	}
	l.emit(TokenKindEOF)
	return nil
}

func lexPlus(l *lexer) stateFunc {
	l.next()
	switch r := l.peek(); {
	case !isAlphaLower(r):
		return l.errorf("expected an identfier")
	}
	l.emit(TokenKindPlus)
	return lexIdent
}

func lexIdent(l *lexer) stateFunc {
	valid := func(r rune) bool {
		return (unicode.IsLetter(r) && unicode.IsLower(r)) || r == colon
	}
	l.acceptFunc(valid)
	r := l.peek()
	if r == '\n' {
		return l.errorf("marker cannot span multiple lines")
	}

	if strings.Count(l.currentValue(), string(colon)) < 1 {
		l.ignore()
		return lexText
	}
	l.emit(TokenKindIdent)

	switch r {
	case '=':
		return lexAssignment
	case eof:
		return lexBool
	default:
		return l.errorf("unexpected token")
	}
}

func lexAssignment(l *lexer) stateFunc {
	l.next()
	l.emit(TokenKindAssignment)
	// figure out which is next
	switch r := l.peek(); {
	case r == '[':
		return lexOpenSquareBracket
	case unicode.IsDigit(r):
		return lexNumber
		// a string is only recognised if a letter
		// follows the assignment symbol. Can change it
		// if an `"` token is introduced
	case unicode.IsLetter(r):
		return lexString
	default:
		return nil
	}
}

func lexBool(l *lexer) stateFunc {
    assignToken := NewToken(TokenKindAssignment, "=")
    l.emitToken(assignToken)
	t := NewToken(TokenKindBool, "true")
	l.emitToken(t)
	return lexText
}

func lexString(l *lexer) stateFunc {
	scanStringWithEscape(l, "", "")
	l.emit(TokenKindString)
	return lexText
}

func lexNumber(l *lexer) stateFunc {
	kind, err := scanNumber(l)
	if err != nil {
		return l.errorf(err.Error())
	}
	l.emit(kind)
	r := l.peek()
	if isNewline(r) {
		l.next()
		// ignore the new line
		l.ignore()
	}
	return lexText
}

func lexOpenSquareBracket(l *lexer) stateFunc {
	l.next()
	l.emit(TokenKindOpenSquareBracket)
	switch r := l.peek(); {
	case r == ']':
		return lexCloseSquareBracket
	case unicode.IsDigit(r):
		return lexNumberArrayValue
	case r == '"':
		return lexStartDoubleQuotationMark
	case isSpace(r):
		return l.errorf("no space allowed of the opening bracket of array")
	case r == ',':
		return l.errorf("expected value in array not seperator")
	default:
		return l.errorf("expected closing bracket. Be sure that there is no whitespace between the last element and `]`")
	}
}

func lexNumberArrayValue(l *lexer) stateFunc {
	kind, err := scanNumber(l)
	if err != nil {
		return l.errorf(err.Error())
	}
	l.emit(kind)
	switch r := l.peek(); {
	case r == ',':
		return lexCommaSeperator
	case r == ']':
		return lexCloseSquareBracket
	default:
		return l.errorf("expected next array value or closing bracket")
	}
}

func lexCommaSeperator(l *lexer) stateFunc {
	l.next()
	l.ignore()
	ignoreSpace(l)
	switch r := l.peek(); {
	case unicode.IsDigit(r):
		return lexNumberArrayValue
	case r == '"':
		return lexStartDoubleQuotationMark
	case r == ']':
		return l.errorf("remove the comma before the closing bracket of the array")
	default:
		return l.errorf("expected next value in array after comma")
	}
}

func lexStartDoubleQuotationMark(l *lexer) stateFunc {
	l.next()
	l.ignore()
	err := scanStringWithEscape(l, `"`, ",]")
	if err != nil {
		return l.errorf(err.Error())
	}
	l.emit(TokenKindString)
	switch r := l.peek(); {
	case r == '"':
		return lexEndDoubleQuotationMark
	default:
		return l.errorf("expected `\"` got `%s`", string(r))
	}
}

func lexEndDoubleQuotationMark(l *lexer) stateFunc {
	l.next()
	l.ignore()
	switch r := l.peek(); {
	case r == ']':
		return lexCloseSquareBracket
	case r == ',':
		return lexCommaSeperator
	default:
		return l.errorf("expected closing bracket or next value of array")
	}
}

func lexCloseSquareBracket(l *lexer) stateFunc {
	l.next()
	l.emit(TokenKindCloseSquareBracket)
	return lexText
}
