package main

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


func Lex(input string) *lexer {
	return &lexer{
		input:  strings.TrimSpace(input),
		tokens: make(chan Token, 100),
		state:  lexText,
		plus:   plus,
	}
}

type lexer struct {
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

func (l *lexer) errorf(format string, args ...any) stateFunc {
	l.tokens <- NewToken(TokenKindError, fmt.Sprintf(format, args...))
	l.start = 0
	l.pos = 0
	l.input = l.input[:0]
	return nil

}

func (l *lexer) next() rune {
	var r rune
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) run() {
	for state := l.state; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) emit(kind TokenKind) {
	l.tokens <- NewToken(kind, l.currentValue())
	l.start = l.pos
}

func (l *lexer) emitToken(t Token) {
	l.tokens <- t
	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// accept is accepting the next rune
// if it's from the `valid` set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) acceptFunc(fn validFunc) {
	for fn(l.next()) {
	}
	l.backup()
}

func (l *lexer) currentValue() string {
	return l.input[l.start:l.pos]
}
