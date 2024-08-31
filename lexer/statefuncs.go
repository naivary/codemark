package lexer

import (
	"strings"
	"unicode"
)

type validFunc func(rune) bool

func ignoreSpace(l *Lexer) {
	l.acceptFunc(isSpace)
	l.ignore()
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r) || r == '-' || r == '+'
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

type stateFunc func(*Lexer) stateFunc

func lexText(l *Lexer) stateFunc {
	// check for the position 0
	if hasPlusPrefix(l.input, l.pos) {
		return lexPlus
	}
	for {
		r := l.next()
		if r == eof {
			break
		}
		if !isNewline(r) {
			continue
		}
		l.acceptFunc(isSpace)
		l.ignore()
		if hasPlusPrefix(l.input, l.pos) {
			return lexPlus
		}
	}
	l.emit(TokenKindEOF)
	return nil
}

func lexPlus(l *Lexer) stateFunc {
	l.next()
	l.emit(TokenKindPlus)
	switch r := l.peek(); {
	case !isAlphaLower(r):
		return l.errorf("after a `+` an immediate identifier is expected. The identifier can only be in lower letters and has to contain two `:` describing the path")
	}
	return lexIdent
}

func lexIdent(l *Lexer) stateFunc {
	valid := func(r rune) bool {
		return (unicode.IsLetter(r) && unicode.IsLower(r)) || r == colon
	}
	l.acceptFunc(valid)
	idn := l.currentValue()
	colons := strings.Count(idn, string(colon))
	if colons < 2 {
		return l.errorf("expected two colons in `%s` but got %d", idn, colons)
	}
	l.emit(TokenKindIdent)

	switch r := l.peek(); {
	case r == '=':
		return lexAssignment
	case r == eof || r == newline:
		return lexBoolWithoutAssignment
	default:
		return l.errorf("expected an assignment operator or a newline after the identifier `%s`", idn)
	}
}

func lexAssignment(l *Lexer) stateFunc {
	l.next()
	l.emit(TokenKindAssignment)
	switch r := l.peek(); {
	case r == '[':
		return lexOpenSquareBracket
	case isDigit(r):
		return lexNumber
	case r == '"':
		return lexStartDoubleQuotationMarkString
	case r == 't' || r == 'f':
		return lexBool
	default:
		return l.errorf("expecting value after assignment. For the possible valid values see: <docs link>")
	}
}

func lexBoolWithoutAssignment(l *Lexer) stateFunc {
	assignToken := NewToken(TokenKindAssignment, "=")
	l.emitToken(assignToken)
	t := NewToken(TokenKindBool, "true")
	l.emitToken(t)
	return lexEndOfExpr
}

func lexBool(l *Lexer) stateFunc {
	spelling := "true"
	if r := l.peek(); r == 'f' {
		spelling = "false"
	}
	for i := 0; i < len(spelling); i++ {
		r := l.peek()
		if r == rune(spelling[i]) {
			l.next()
			continue
		}
		return l.errorf("`%s` is not spelled correctly", spelling)
	}
	l.emit(TokenKindBool)
	return lexEndOfExpr
}

func lexNumber(l *Lexer) stateFunc {
	kind, err := scanNumber(l)
	if err != nil {
		return l.errorf(err.Error())
	}
	l.emit(kind)
	return lexEndOfExpr
}

func lexOpenSquareBracket(l *Lexer) stateFunc {
	l.next()
	l.emit(TokenKindOpenSquareBracket)
	switch r := l.peek(); {
	case r == ']':
		return lexCloseSquareBracket
	case isDigit(r):
		return lexNumberArrayValue
	case r == '"':
		return lexStartDoubleQuotationMarkStringArray
	case r == 't' || r == 'f':
		return lexBoolArrayValue
	case isSpace(r):
		return l.errorf("no space allowed after the opening bracket of an array")
	case r == ',':
		return l.errorf("expected value in array not seperator")
	default:
		return l.errorf("expected closing bracket. Be sure that there is no whitespace between the last element and `]`")
	}
}

func lexArraySequence(l *Lexer) stateFunc {
	switch r := l.peek(); {
	case r == ',':
		return lexCommaSeperator
	case r == ']':
		return lexCloseSquareBracket
	default:
		return l.errorf("expected next array value or closing bracket")
	}
}

func lexBoolArrayValue(l *Lexer) stateFunc {
	spelling := "true"
	if r := l.peek(); r == 'f' {
		spelling = "false"
	}
	for i := 0; i < len(spelling); i++ {
		r := l.peek()
		if r == rune(spelling[i]) {
			l.next()
			continue
		}
		return l.errorf("`%s` is not spelled correctly", spelling)
	}
	l.emit(TokenKindBool)
	return lexArraySequence
}

func lexNumberArrayValue(l *Lexer) stateFunc {
	kind, err := scanNumber(l)
	if err != nil {
		return l.errorf(err.Error())
	}
	l.emit(kind)
	return lexArraySequence
}

func lexCommaSeperator(l *Lexer) stateFunc {
	l.next()
	l.ignore()
	ignoreSpace(l)
	switch r := l.peek(); {
	case isDigit(r):
		return lexNumberArrayValue
	case r == '"':
		return lexStartDoubleQuotationMarkStringArray
	case r == ']':
		return l.errorf("remove the comma before the closing bracket of the array")
	case r == 't' || r == 'f':
		return lexBoolArrayValue
	default:
		return l.errorf("expected next value in array after comma")
	}
}

func lexStartDoubleQuotationMarkString(l *Lexer) stateFunc {
	l.next()
	l.ignore()
	if err := scanString(l); err != nil {
		return l.errorf(err.Error())
	}
	l.emit(TokenKindString)
	switch r := l.peek(); {
	case r == '"':
		return lexEndDoubleQuotationMarkString
	default:
		return l.errorf("expected `\"` got `%s`", string(r))
	}
}

func lexEndDoubleQuotationMarkString(l *Lexer) stateFunc {
	l.next()
	l.ignore()
	return lexEndOfExpr
}

func lexStartDoubleQuotationMarkStringArray(l *Lexer) stateFunc {
	l.next()
	l.ignore()
	if err := scanString(l); err != nil {
		return l.errorf(err.Error())
	}
	l.emit(TokenKindString)
	switch r := l.peek(); {
	case r == '"':
		return lexEndDoubleQuotationMarkStringArray
	default:
		return l.errorf("expected `\"` got `%s`", string(r))
	}
}

func lexEndDoubleQuotationMarkStringArray(l *Lexer) stateFunc {
	l.next()
	l.ignore()
	return lexArraySequence
}

func lexCloseSquareBracket(l *Lexer) stateFunc {
	l.next()
	l.emit(TokenKindCloseSquareBracket)
	return lexEndOfExpr
}

func lexEndOfExpr(l *Lexer) stateFunc {
	l.acceptFunc(isSpace)
	l.ignore()
	switch r := l.peek(); {
	case r == newline || r == eof:
		return lexText
	default:
		return l.errorf("after a finished marker expression only a newline can follow")
	}
}
