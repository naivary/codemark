package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/naivary/codemark/lexer/token"
)

// TODO: should we allow multi line string in list?

const (
	_eof        = -1
	_plus       = '+'
	_colon      = ':'
	_newline    = '\n'
	_underscore = '_'
	_dot        = '.'
	_whitespace = ' '
	_tab        = '\t'
	_tick       = '`'
	_dquot      = '"'
	_lbrack     = '['
	_rbrack     = ']'
	_assign     = '='
	_comma      = ','
)

func Lex(input string) *Lexer {
	l := &Lexer{
		input:  strings.TrimSpace(input),
		tokens: make(chan Token, 100),
		state:  lexText,
	}
	l.run()
	return l
}

type Lexer struct {
	// the string being scanned
	input string
	// start position of this item
	start int
	// current position of this item
	pos int
	// length of the last rune read
	width int
	// the `stateFunc` to begin with
	state stateFunc
	// channel of scanned tokens
	tokens chan Token
}

// NextToken returns the next token of the channel
func (l *Lexer) NextToken() Token {
	for {
		select {
		case token := <-l.tokens:
			return token
		default:
			l.state = l.state(l)
		}
	}
}

func (l *Lexer) run() {
	for state := l.state; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *Lexer) errorf(format string, args ...any) stateFunc {
	l.tokens <- NewToken(token.ERROR, fmt.Sprintf(format, args...))
	l.start = 0
	l.pos = 0
	l.input = l.input[:0]
	return nil

}

func (l *Lexer) next() rune {
	var r rune
	if l.pos >= len(l.input) {
		l.width = 0
		return _eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *Lexer) emit(kind token.Kind) {
	l.tokens <- NewToken(kind, l.currentValue())
	l.start = l.pos
}

func (l *Lexer) emitToken(t Token) {
	l.tokens <- t
	l.start = l.pos
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// accept is accepting the next rune
// if it's from the `valid` set.
func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *Lexer) acceptFunc(fn func(r rune) bool) {
	for fn(l.next()) {
	}
	l.backup()
}

func (l *Lexer) currentValue() string {
	return l.input[l.start:l.pos]
}
