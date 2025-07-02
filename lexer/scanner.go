package lexer

import (
	"fmt"
	"strconv"

	"github.com/naivary/codemark/lexer/token"
)

// scanRealNumber is readin in the next characters from the lexer and trying to
// figure out if it is a FLOAT or INT. It is only a helper function used by the
// parent function `scanNumber`.
func scanRealNumber(l *Lexer) (token.Kind, error) {
	kind := token.INT
	l.accept("+-")
	digits := "0123456789_"
	// Is it hex?
	if l.accept("0") {
		// NOTE: Leading 0 does not mean octal in floats.
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

// scanNumber is reading in the next characters from the lexer and is trying to
// figure out which kind of number is provided e.g. INT, FLOAT, COMPLEX. If no
// number can be recognised then token.ERROR will be returned.
func scanNumber(l *Lexer) (token.Kind, error) {
	var err error
	var kind token.Kind
	kind, err = scanRealNumber(l)
	if err != nil {
		return token.ERROR, err
	}
	r := l.peek()
	if r != 'i' && r != '+' && r != '-' {
		return kind, nil
	}
	if l.accept("i") {
		_, err = scanRealNumber(l)
		return token.COMPLEX, err
	}
	_, err = scanRealNumber(l)
	if !l.accept("i") {
		return token.ERROR, fmt.Errorf("two numbers can only be defined if you want to define an complex number e.g. 2+3i: %s\n", l.currentValue())
	}
	return token.COMPLEX, nil
}

func scanMultiLineString(l *Lexer) error {
	v := func(r rune) bool {
		return r != _tick && r != _eof
	}
	l.acceptFunc(v)
	quoted := fmt.Sprintf("`%s`", l.currentValue())
	if _, err := strconv.Unquote(quoted); err != nil {
		return fmt.Errorf("string `%s` is not a valid go string", l.currentValue())
	}
	return nil
}

func scanString(l *Lexer) error {
	const backslash = '\\'
	v := func(r rune) bool {
		if r != '"' && r != backslash && r != _eof && !isNewline(r) {
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
