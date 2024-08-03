package main

import (
	"errors"
	"fmt"
	"strings"
)

func scanRealNumber(l *lexer) (TokenKind, error) {
	kind := TokenKindInt
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
		kind = TokenKindFloat
		l.acceptRun(digits)
	}
	if len(digits) == 10+1 && l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	return kind, nil
}

func scanNumber(l *lexer) (TokenKind, error) {
	kind := TokenKindInt
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
		kind = TokenKindFloat
		l.acceptRun(digits)
	}
	if len(digits) == 10+1 && l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	// Is it imaginary?
	if l.accept("i") {
		kind = TokenKindComplex
	}
	r := l.peek()
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(r) {
		l.next()
		return TokenKindError, errors.New("bad syntax for number")
	}
	return kind, nil
}

func isSpecialCharacter(escape string, r rune) bool {
	return strings.Contains(escape, string(r))
}

// `c` is defining which characters have to follow after a
// a character defined by `escape` so it is still a valid unescaped symbol.
// For example +path:to:marker=["item\"s"]. The last `"` don't have to be
// escaped because its the end of the string followed by a `]`.
func scanStringWithEscape(l *lexer, escape string, c string) error {
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
	return fmt.Errorf("special character `%s` is not escaped", string(r))
}

func isCorrectUnescaped(l *lexer, c string) bool {
	r := l.peek()
	return strings.Contains(c, string(r))
}

func escapeChar(l *lexer, escape string, c string) error {
	if isSpecialCharacter(escape, l.peek()) {
		l.next()
		return scanStringWithEscape(l, escape, c)
	}
	return errors.New("unexpected `\\` in string literal. Has to be escaped using `\\`")
}

func scanString(l *lexer) {
	valid := func(r rune) bool {
		return !isSpace(r) && r != eof && !isNewline(r)
	}
	l.acceptFunc(valid)
}
