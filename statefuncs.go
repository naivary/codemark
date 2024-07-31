package main

import (
	"strings"
	"unicode"
)

func isSpace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\r' || r == '\t'
}

func isNewLine(r rune) bool {
	return r == '\n' || r == '\r'
}

func hasPlusPrefix(input string, pos int) bool {
	return strings.HasPrefix(input[pos:], plus)
}

// isAlphaNumeric is checking if the rune is a letter, digit or underscore
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func lexText(l *lexer) stateFunc {
	for {
		// check for the position 0
		if hasPlusPrefix(l.input, l.pos) {
			return lexPlus
		}

		r := l.next()
		if r == eof {
			break
		}
		// ignore anything before the new position
		// because it's not relevant
		l.ignore()
		if hasPlusPrefix(l.input, l.pos) {
			return lexPlus
		}
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
	valid := func(r rune) bool {
		return (unicode.IsLetter(r) && unicode.IsLower(r)) || r == colon
	}
	l.acceptFunc(valid)
	r := l.peek()
	if r == '\n' {
		return l.errorf("marker cannot span multiple lines")
	}

	if strings.Count(l.currentValue(), string(colon)) != 2 {
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
	case r == '{':
		return lexOpenCurlyBracket
	case unicode.IsDigit(r):
		return lexNumber
		// a string is only recognised if a letter
		// follows the assignment symbol
	case unicode.IsLetter(r):
		return lexString
	default:
		return nil
	}
}

func lexBool(l *lexer) stateFunc {
	t := NewToken(TokenKindBool, "true")
	l.emitToken(t)
	return lexText
}

func lexString(l *lexer) stateFunc {
	valid := func(r rune) bool {
		return !isSpace(r) && r != eof
	}
	l.acceptFunc(valid)
	l.emit(TokenKindString)
	return lexText
}

func lexNumber(l *lexer) stateFunc {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789_"
	if l.accept("0") {
		// Note: Leading 0 does not mean octal in floats.
		if l.accept("xX") {
			digits = "0123456789abcdefABCDEF_"
		} else if l.accept("oO") {
			digits = "01234567_"
		} else if l.accept("bB") {
			digits = "01_"
		}
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if len(digits) == 10+1 && l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	if len(digits) == 16+6+1 && l.accept("pP") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	// Is it imaginary?
	l.accept("i")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad syntax for number")
	}
	l.emit(TokenKindNumber)
	return lexText
}

func lexOpenCurlyBracket(l *lexer) stateFunc {
	return nil
}

func lexOpenSquareBracket(l *lexer) stateFunc {
	return nil
}
