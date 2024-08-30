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

func ignoreNewline(l *Lexer) {
	l.acceptFunc(isNewline)
	l.ignore()
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
		if isNewline(r) {
			l.acceptFunc(isSpace)
			l.ignore()
			if hasPlusPrefix(l.input, l.pos) {
				return lexPlus
			}
		}
	}
	l.emit(TokenKindEOF)
	return nil
}

func lexPlus(l *Lexer) stateFunc {
	l.next()
	switch r := l.peek(); {
	case !isAlphaLower(r):
		return l.errorf("expected an identfier")
	}
	l.emit(TokenKindPlus)
	return lexIdent
}

func lexIdent(l *Lexer) stateFunc {
	valid := func(r rune) bool {
		return (unicode.IsLetter(r) && unicode.IsLower(r)) || r == colon
	}
	l.acceptFunc(valid)
	r := l.peek()
	if r == '\n' {
		return l.errorf("marker cannot span multiple lines")
	}

	if strings.Count(l.currentValue(), string(colon)) < 1 {
		l.ignore()
		return lexText
	}
	l.emit(TokenKindIdent)

	switch r {
	case '=':
		return lexAssignment
	case eof:
		return lexBoolWithoutAssignment
	default:
		return l.errorf("unexpected token")
	}
}

func lexAssignment(l *Lexer) stateFunc {
	l.next()
	l.emit(TokenKindAssignment)
	switch r := l.peek(); {
	case r == '[':
		return lexOpenSquareBracket
	case unicode.IsDigit(r):
		return lexNumber
	case r == '"':
		return lexStartDoubleQuotationMarkString
	case r == 't' || r == 'f':
		return lexBool
	case r == '-' || r == '+':
		return lexNumber
	default:
		return l.errorf("expecting value after assignment")
	}
}

func lexBoolWithoutAssignment(l *Lexer) stateFunc {
	assignToken := NewToken(TokenKindAssignment, "=")
	l.emitToken(assignToken)
	t := NewToken(TokenKindBool, "true")
	l.emitToken(t)
	return lexText
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
	l.acceptFunc(isSpace)
	r := l.peek()
	if r == eof {
		l.emit(TokenKindBool)
		return lexText
	}
	return l.errorf("did not expect anything after `%s`. Found `%s`", spelling, string(r))
}

func lexNumber(l *Lexer) stateFunc {
	kind, err := scanNumber(l)
	if err != nil {
		return l.errorf(err.Error())
	}
	l.emit(kind)
	r := l.peek()
	if isNewline(r) {
		l.next()
		// ignore the new line
		l.ignore()
	}
	return lexText
}

func lexOpenSquareBracket(l *Lexer) stateFunc {
	l.next()
	l.emit(TokenKindOpenSquareBracket)
	switch r := l.peek(); {
	case r == ']':
		return lexCloseSquareBracket
	case unicode.IsDigit(r) || r == '-' || r == '+':
		return lexNumberArrayValue
	case r == '"':
		return lexStartDoubleQuotationMarkStringArray
	case r == 't' || r == 'f':
		return lexBoolArrayValue
	case isSpace(r):
		return l.errorf("no space allowed of the opening bracket of array")
	case r == ',':
		return l.errorf("expected value in array not seperator")
	default:
		return l.errorf("expected closing bracket. Be sure that there is no whitespace between the last element and `]`")
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
	switch r := l.peek(); {
	case r == ',':
		return lexCommaSeperator
	case r == ']':
		return lexCloseSquareBracket
	default:
		return l.errorf("expected next array value or closing bracket")
	}
}

func lexNumberArrayValue(l *Lexer) stateFunc {
	kind, err := scanNumber(l)
	if err != nil {
		return l.errorf(err.Error())
	}
	l.emit(kind)
	switch r := l.peek(); {
	case r == ',':
		return lexCommaSeperator
	case r == ']':
		return lexCloseSquareBracket
	default:
		return l.errorf("expected next array value or closing bracket")
	}
}

func lexCommaSeperator(l *Lexer) stateFunc {
	l.next()
	l.ignore()
	ignoreSpace(l)
	switch r := l.peek(); {
	case unicode.IsDigit(r) || r == '-' || r == '+':
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
	switch r := l.peek(); {
	case r == ']':
		return lexCloseSquareBracket
	case r == ',':
		return lexCommaSeperator
	default:
		return l.errorf("expected closing bracket or next value of array")
	}
}

func lexCloseSquareBracket(l *Lexer) stateFunc {
	l.next()
	l.emit(TokenKindCloseSquareBracket)
	return lexText
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
