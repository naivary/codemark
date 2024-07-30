package main

import (
	"strings"
	"unicode"
)

func isSpace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\r' || r == '\t'
}

// isIdenfitier is returning if `r` is a correct
// identifier character
func isIdenfitier(r rune) bool {
	if unicode.IsLetter(r) && unicode.IsLower(r) {
		return true
	}
	return r == ':'
}

func hasPlusPrefix(input string, pos int) bool {
	return strings.HasPrefix(input[pos:], "+")
}

func lexText(l *lexer) stateFunc {
	if hasPlusPrefix(l.input, l.pos) {
		return lexPlus
	}
	for r := l.next(); r != eof; l.next() {
		if hasPlusPrefix(l.input, l.pos) {
			l.ignore()
			return lexPlus
		}
		l.ignore()
	}
	l.emit(TokenKindEOF)
	return nil
}

func lexPlus(l *lexer) stateFunc {
	l.pos += len(l.plus)
	l.emit(TokenKindPlus)
	return lexIdent
}

func lexIdent(l *lexer) stateFunc {
	switch r := l.next(); {
	case r == '\n':
		return l.errorf("marker cannot span multiple lines")
	case isSpace(r):
		return l.errorf("after a plus no space is allowed")
	case unicode.IsDigit(r):
		return l.errorf("no digits allowed as identifier")
	case isIdenfitier(r):
		return lexIdent
	case r == '=':
		return lexAssignment
	case r == eof:
		return lexBool
	default:
		return l.errorf("unrecognized action")
	}
}

func lexAssignment(l *lexer) stateFunc {
	l.backup()
	l.emit(TokenKindIdent)
	l.forward()
	l.emit(TokenKindAssignment)
	// figure out which is next
	switch r := l.next(); {
	case r == '[':
		return lexOpenSquareBracket
	case r == '{':
		return lexMap
	case unicode.IsDigit(r):
		return lexNumber
	case unicode.IsLetter(r):
		return lexString
	default:
		return nil
	}
}

// [ value1, ]
func lexOpenSquareBracket(l *lexer) stateFunc {
	l.emit(TokenKindOpenSquareBracket)
	return lexArrayValue
}

func lexArrayValue(l *lexer) stateFunc {
	switch r := l.next(); {
	case isSpace(r):
		l.ignore()
		return lexArrayValue
	case r == ',':
		return lexCommaSep
	default:
		return lexArrayValue
	}
}

func lexCommaSep(l *lexer) stateFunc {
	l.backup()
	l.emit(TokenKindArrayValue)
	l.forward()
	l.emit(TokenKindCommaSeparator)
	switch r := l.next(); {
	case isSpace(r):
		l.ignore()
		return lexCommaSep
	case r == ']':
		return lexCloseSquareBracket
	case unicode.IsLetter(r):
		return lexArrayValue
	default:
		return l.errorf("expected value or closing squared bracket")
	}
}

func lexCloseSquareBracket(l *lexer) stateFunc {
	l.emit(TokenKindArrayValue)
	return lexText
}

func lexMap(l *lexer) stateFunc {
	return nil
}

func lexNumber(l *lexer) stateFunc {
	return nil
}

func lexString(l *lexer) stateFunc {
	return nil
}

func lexBool(l *lexer) stateFunc {
	l.emit(TokenKindIdent)
	l.emitToken(NewToken(TokenKindBool, "true"))
	return lexText
}
