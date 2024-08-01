package main

import (
	"strings"
)

func scanNumber(l *lexer) stateFunc {
	// Optional leading sign.
	l.accept("+-")
	digits := "0123456789_"
	// Is it hex?
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
	r := l.peek()
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(r) {
		l.next()
		return l.errorf("bad syntax for number")
	}
	return nil
}

func isSpecialCharacter(escape string, r rune) bool {
	return strings.Contains(escape, string(r))
}

func scanStringWithEscape(l *lexer, escape string, c string) stateFunc {
	// if no escape characters are provided
	// none will be escaped
	if escape != "" && !strings.Contains(escape, "\\") {
		// `\` has to be escaped too
		escape += "\\"
	}
	valid := func(r rune) bool {
		return !isSpace(r) && r != eof && !isSpecialCharacter(escape, r)
	}
	l.acceptFunc(valid)

	r := l.peek()
	if r == eof || isSpace(r) {
		// scanned the full string without any error
		return nil
	}
	if r == '\\' {
		// check if correct escaping is taking place
		l.next()
		return escapeChar(l, escape, c)
	}
	l.next()
	isCorrect := isCorrectUnescaped(l, c)
	if isCorrect {
        l.backup()
		return nil
	}
	return l.errorf("special character `%s` is not escaped", string(r))
}

func isCorrectUnescaped(l *lexer, c string) bool {
	r := l.peek()
	return strings.Contains(c, string(r))
}

func escapeChar(l *lexer, escape string, c string) stateFunc {
	if isSpecialCharacter(escape, l.peek()) {
		l.next()
		return scanStringWithEscape(l, escape, c)
	}
	return l.errorf("unexpected `\\` in string literal. Has to be escaped using `\\`")
}

func scanString(l *lexer) {
	valid := func(r rune) bool {
		return !isSpace(r) && r != eof && !isNewline(r)
	}
	l.acceptFunc(valid)
}
