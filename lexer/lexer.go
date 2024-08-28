package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	eof     = -1
	plus    = "+"
	colon   = ':'
	newline = '\n'
)

func Lex(input string) *Lexer {
	return &Lexer{
		input:  strings.TrimSpace(input),
		tokens: make(chan Token, 100),
		state:  lexText,
		plus:   plus,
	}
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
	// state is the `stateFunc` to begin with
	state stateFunc
	// channel of scanned tokens
	tokens chan Token
	// needed sign to satrt the action
	plus string
	// token to return to the parser
	token Token
}

func (l *Lexer) errorf(format string, args ...any) stateFunc {
	l.tokens <- NewToken(TokenKindError, fmt.Sprintf(format, args...))
	l.start = 0
	l.pos = 0
	l.input = l.input[:0]
	return nil

}

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

func (l *Lexer) next() rune {
	var r rune
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *Lexer) Run() {
	for state := l.state; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *Lexer) emit(kind TokenKind) {
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

func (l *Lexer) forceBackup() {
	_, width := utf8.DecodeRuneInString(string(l.input[len(l.input)-1]))
	l.pos -= width
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

func (l *Lexer) acceptFunc(fn validFunc) {
	for fn(l.next()) {
	}
	l.backup()
}

func (l *Lexer) currentValue() string {
	return l.input[l.start:l.pos]
}
