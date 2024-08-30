package lexer

import (
	"errors"
	"fmt"
	"strconv"
)

func scanRealNumber(l *Lexer) (TokenKind, error) {
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

func scanNumber(l *Lexer) (TokenKind, error) {
	var err error
	kind := TokenKindInt
	kind, err = scanRealNumber(l)
	if err != nil {
		return TokenKindError, err
	}
	if l.accept("i") {
		if l.accept("+-") {
			return TokenKindError, errors.New("real part of a complex number has to be defined before the imaginary part e.g. `3+2i` not `2i+3`")
		}
		kind = TokenKindComplex
	}
	r := l.peek()
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(r) {
		l.next()
		return TokenKindError, errors.New("bad syntax for number")
	}
	if l.accept("+-") {
		_, err = scanRealNumber(l)
		if !l.accept("i") {
			return TokenKindError, errors.New("missing imaginary symbol `i`")
		}
		kind = TokenKindComplex
	}
	return kind, err
}

func scanString(l *Lexer) error {
	const backslash = '\\'
	v := func(r rune) bool {
		if r == backslash {
			l.next()
			return true
		}
		if r != eof && r != '"' {
			return true
		}
		return false
	}
	l.acceptFunc(v)
	quoted := fmt.Sprintf(`"%s"`, l.currentValue())
	if _, err := strconv.Unquote(quoted); err != nil {
		return fmt.Errorf("string `%s` is not a valid go string", l.currentValue())
	}
	return nil
}
