package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const eof = -1

type stateFunc func(*lexer) stateFunc

func lex(name, input string) *lexer {
	l := &lexer{
		name:   name,
		input:  strings.TrimSpace(input),
		tokens: make(chan Token, 100),
		state:  lexText,
		plus:   "+",
	}
	return l
}

type lexer struct {
	// used only for error reports
	name string
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

	plus string

	// token to return to the parser
	token Token
}

func (l *lexer) errorf(format string, args ...any) stateFunc {
	l.token = NewToken(TokenKindError, fmt.Sprintf(format, args...))
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
	l.tokens <- NewToken(kind, l.input[l.start:l.pos])
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

func (l *lexer) forward() {
	l.pos += l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}
