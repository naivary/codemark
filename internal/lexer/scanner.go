package lexer

import (
	"fmt"
	"strconv"
)

// TODO: scanString sollte nur eine Funktino sein am besten

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
