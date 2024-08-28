package lexer

import (
	"errors"
	"fmt"
	"strings"
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

func isSpecialCharacter(escape string, r rune) bool {
	return strings.Contains(escape, string(r))
}

func scanStringN(l *Lexer) error {
	chars := `"\`
    // +path:to:marker="string\""
    // 
	valid := func(r rune) bool {
		return r != eof && !strings.Contains(chars, string(r))
	}
	l.acceptFunc(valid)

	// r will be either " or \
	// if its a " return it
	if r := l.peek(); r == '"' {
		return nil
	}
	// if here then its a `\`
	// we have to check which character follows the \
	l.next()

	return nil
}

// `c` is defining which characters have to follow after a
// a character defined by `escape` so it is still a valid unescaped symbol.
// For example +path:to:marker=["item\"s"]. The last `"` don't have to be
// escaped because its the end of the string followed by a `]`.
func scanStringWithEscape(l *Lexer, escape string, c string) error {
	// if no escape characters are provided
	// none will be escaped
	if escape != "" && !strings.Contains(escape, "\\") {
		// `\` has to be escaped too
		escape += "\\"
	}
	valid := func(r rune) bool {
		return r != eof && !isSpecialCharacter(escape, r) && !isNewline(r)
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
	// have to keep the width because the next character may be a eof which
	// results in the width lost.
	width := l.width
	isCorrect := isCorrectUnescaped(l, c)
	if isCorrect {
		l.width = width
		l.backup()
		return nil
	}
	return fmt.Errorf("special character `%s` is not escaped", string(r))
}

func isCorrectUnescaped(l *Lexer, c string) bool {
	r := l.peek()
	if c == "" && (r == eof || isSpace(r)) {
		return true
	}
	return strings.Contains(c, string(r))
}

func escapeChar(l *Lexer, escape string, c string) error {
	if isSpecialCharacter(escape, l.peek()) {
		l.next()
		return scanStringWithEscape(l, escape, c)
	}
	return errors.New("unexpected `\\` in string literal. Has to be escaped using `\\`")
}

func scanString(l *Lexer) {
	valid := func(r rune) bool {
		return !isSpace(r) && r != eof && !isNewline(r)
	}
	l.acceptFunc(valid)
}
