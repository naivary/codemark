package lexer

import (
	"fmt"
	"strconv"

	"github.com/naivary/codemark/lexer/token"
)

func scanRealNumber(l *Lexer) (token.Kind, error) {
	kind := token.INT
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
		kind = token.FLOAT
		l.acceptRun(digits)
	}
	if len(digits) == 10+1 && l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	return kind, nil
}

func scanNumber(l *Lexer) (token.Kind, error) {
	// TODO: I think we can set before the return kind = token.COMPLEX
	var err error
	kind := token.INT
	kind, err = scanRealNumber(l)
	if err != nil {
		return token.ERROR, err
	}
	if l.accept("i") {
		if l.accept("+-") {
			return token.ERROR, ErrRealBeforeComplex
		}
		kind = token.COMPLEX
	}
	r := l.peek()
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(r) {
		l.next()
		return token.ERROR, ErrBadSyntaxForNumber
	}
	if l.accept("+-") {
		_, err = scanRealNumber(l)
		if !l.accept("i") {
			return token.ERROR, ErrImagMissing
		}
		kind = token.COMPLEX
	}
	return kind, err
}

func scanString(l *Lexer) error {
	const backslash = '\\'
	v := func(r rune) bool {
		if r != '"' && r != backslash && r != _eof {
			return true
		}
		if r == backslash {
			l.next()
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
