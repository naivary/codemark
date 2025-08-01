package lexer

import (
	"unicode"

	"github.com/naivary/codemark/internal/lexer/token"
	"github.com/naivary/codemark/validate"
)

var isListSeq = false

type stateFunc func(*Lexer) stateFunc

func lexText(l *Lexer) stateFunc {
	// check for the position 0
	if hasPlusPrefix(l.input, l.pos) {
		return lexPlus
	}
	for {
		r := l.next()
		if r == _eof {
			return lexEOF
		}
		if !isNewline(r) {
			continue
		}
		// accept any amount of whitespace
		l.acceptFunc(isSpace)
		l.ignore()
		if hasPlusPrefix(l.input, l.pos) {
			return lexPlus
		}
	}
}

func lexPlus(l *Lexer) stateFunc {
	l.next()
	l.emit(token.PLUS)
	switch r := l.peek(); {
	case !isLower(r):
		return l.errorf(
			"after a `%s` an immediate identifier is expected. The identifier can only be in lower letters and has to contain two `:` describing the path",
			string(_plus),
		)
	}
	return lexIdent
}

func lexIdent(l *Lexer) stateFunc {
	if r := l.peek(); isDigit(r) {
		return l.errorf("marker identifier cannot start with a digit: %d", r)
	}
	valid := func(r rune) bool {
		return (unicode.IsLetter(r) && unicode.IsLower(r)) || unicode.IsDigit(r) || r == _colon ||
			r == _underscore ||
			r == _dot
	}
	l.acceptFunc(valid)
	ident := l.currentValue()
	if err := validate.Ident(ident); err != nil {
		return l.errorf("err: %s", err.Error())
	}
	l.emit(token.IDENT)
	switch l.peek() {
	case _assign:
		return lexAssign
	case _eof, _newline:
		return lexBoolWithoutAssignment
	default:
		return l.errorf("expected an assignment operator or a newline after the identifier `%s`", ident)
	}
}

func lexAssign(l *Lexer) stateFunc {
	l.next()
	l.emit(token.ASSIGN)
	switch r := l.peek(); {
	case r == _lbrack:
		return lexLBRACK
	case isDigit(r):
		return lexNumber
	case r == _dquot:
		return lexStartDQUOT
	case r == _tick:
		return lexStartTick
	case isBool(r):
		return lexBool
	default:
		return l.errorf(
			"expecting value after assignment. For the possible valid values see: <docs link>",
		)
	}
}

func lexStartTick(l *Lexer) stateFunc {
	l.next()
	l.ignore()
	if err := scanMultiLineString(l); err != nil {
		return l.errorf("error: %s", err.Error())
	}
	l.emit(token.STRING)
	switch l.peek() {
	case _tick:
		return lexEndTick
	default:
		return l.errorf("expected `\"` got `%s`", string(l.peek()))
	}
}

func lexEndTick(l *Lexer) stateFunc {
	l.next()
	l.ignore()
	return lexEndOfExpr
}

func lexBoolWithoutAssignment(l *Lexer) stateFunc {
	assignToken := NewToken(token.ASSIGN, "=")
	l.emitToken(assignToken)
	t := NewToken(token.BOOL, "true")
	l.emitToken(t)
	return lexEndOfExpr
}

func lexBool(l *Lexer) stateFunc {
	spelling := "true"
	if r := l.peek(); r == 'f' {
		spelling = "false"
	}
	for i := range len(spelling) {
		r := l.peek()
		if r == rune(spelling[i]) {
			l.next()
			continue
		}
		return l.errorf("`%s` is not spelled correctly", spelling)
	}
	l.emit(token.BOOL)
	if isListSeq {
		return lexListSeq
	}
	return lexEndOfExpr
}

func lexNumber(l *Lexer) stateFunc {
	l.acceptFunc(func(r rune) bool {
		// _comma is needed for list purposes
		return !isNewline(r) && r != _eof && r != _comma && r != _rbrack
	})
	number := l.currentValue()
	kind, err := kindOfNumber(number)
	if err != nil {
		return l.errorf("%s", err.Error())
	}
	l.emit(kind)
	if isListSeq {
		return lexListSeq
	}
	return lexEndOfExpr
}

func lexLBRACK(l *Lexer) stateFunc {
	l.next()
	l.emit(token.LBRACK)
	isListSeq = true
	switch r := l.peek(); {
	case r == _rbrack:
		return lexRBRACK
	case isDigit(r):
		return lexNumber
	case r == _dquot:
		return lexStartDQUOT
	case isBool(r):
		return lexBool
	case isSpace(r):
		return l.errorf("no space allowed after the opening bracket of a list")
	case r == _comma:
		return l.errorf("expected value in array not seperator")
	case r == _tick:
		return l.errorf("multiline strings are not supported in list")
	default:
		return l.errorf(
			"expected closing bracket. Be sure that there is no whitespace between the last element and `]`",
		)
	}
}

func lexListSeq(l *Lexer) stateFunc {
	switch l.peek() {
	case _comma:
		return lexComma
	case _rbrack:
		return lexRBRACK
	default:
		return l.errorf("expected next array value or closing bracket")
	}
}

func lexComma(l *Lexer) stateFunc {
	l.next()
	l.ignore()
	ignoreSpace(l)
	switch r := l.peek(); {
	case isDigit(r):
		return lexNumber
	case r == _dquot:
		return lexStartDQUOT
	case r == _rbrack:
		return l.errorf("remove the comma before the closing bracket of the list")
	case isBool(r):
		return lexBool
	case r == _tick:
		return l.errorf("multiline strings are not supported in list")
	default:
		return l.errorf("expected next value in list after comma")
	}
}

func lexStartDQUOT(l *Lexer) stateFunc {
	l.next()
	l.ignore()
	if err := scanString(l); err != nil {
		return l.errorf("error: %s", err.Error())
	}
	l.emit(token.STRING)
	switch l.peek() {
	case _dquot:
		return lexEndDQUOT
	default:
		return l.errorf("expected `\"` got `%s`", string(l.peek()))
	}
}

func lexEndDQUOT(l *Lexer) stateFunc {
	l.next()
	l.ignore()
	if isListSeq {
		return lexListSeq
	}
	return lexEndOfExpr
}

func lexRBRACK(l *Lexer) stateFunc {
	l.next()
	l.emit(token.RBRACK)
	isListSeq = false
	return lexEndOfExpr
}

func lexEndOfExpr(l *Lexer) stateFunc {
	l.acceptFunc(isSpace)
	l.ignore()
	switch r := l.peek(); {
	case isNewline(r):
		return lexText
	case r == _eof:
		return lexEOF
	default:
		return l.errorf(
			"after a finished marker expression only a newline can follow but found: `%s`",
			string(r),
		)
	}
}

func lexEOF(l *Lexer) stateFunc {
	l.emit(token.EOF)
	return nil
}
