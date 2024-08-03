package main

import (
	"reflect"
	"strconv"
)

type parseFunc func(*parser, Token) (parseFunc, Token)

type Marker interface {
	// Ident is the identifier of
	// the marker without the `+`
	Ident() string

	// Kind is the reflect.Kind
	// the marker is using
	Kind() reflect.Kind

	//
	Value() reflect.Value
}

var _ Marker = (*marker)(nil)

type marker struct {
	ident string
	kind  reflect.Kind
	value reflect.Value
}

func (m *marker) Ident() string {
	return m.ident
}

func (m *marker) Kind() reflect.Kind {
	return m.kind
}

func (m *marker) Value() reflect.Value {
	return m.value
}

func parse(input string) *parser {
	return &parser{
		l:       Lex(input),
		state:   parsePlus,
		markers: make(chan Marker, 2),
		m:       &marker{},
	}
}

type parser struct {
	l       *lexer
	state   parseFunc
	markers chan Marker
	m       *marker
}

func (p *parser) run() {
	p.l.run()
	t := p.next()
	for state := p.state; state != nil; {
		state, t = state(p, t)
	}
	close(p.markers)
}

func (p *parser) emit() {
	p.markers <- p.m
	p.m = &marker{}
}

func (p *parser) next() Token {
	return <-p.l.tokens
}

func parsePlus(p *parser, t Token) (parseFunc, Token) {
	switch t.Kind {
	case TokenKindEOF:
		return nil, t
	default:
		return parseIdent, p.next()
	}
}

func parseIdent(p *parser, t Token) (parseFunc, Token) {
	p.m.ident = t.Value
	return parseAssignment, p.next()
}

func parseAssignment(p *parser, t Token) (parseFunc, Token) {
	return parseValue, p.next()
}

func parseValue(p *parser, t Token) (parseFunc, Token) {
	switch t.Kind {
	case TokenKindString:
		return parseString, t
	case TokenKindInt:
		return parseInt, t
    case TokenKindComplex:
        return parseComplex, t
	case TokenKindOpenSquareBracket:
		return parseSlice, t
	case TokenKindBool:
		return parseBool, t
	default:
		// Something went wrong while lexing
		return nil, t
	}
}

func parseBool(p *parser, t Token) (parseFunc, Token) {
	p.m.kind = reflect.Bool
	boolVal, err := strconv.ParseBool(t.Value)
	if err != nil {
		// error handling
		return nil, t
	}
	p.m.value = reflect.ValueOf(boolVal)
	return parseEOF, p.next()
}

func parseString(p *parser, t Token) (parseFunc, Token) {
    p.m.kind = reflect.String
    p.m.value = reflect.ValueOf(t.Value)
	return parseEOF, p.next()
}

func parseInt(p *parser, t Token) (parseFunc, Token) {
    i, err := strconv.ParseInt(t.Value, 0, 64)
    if err != nil {
        // error handling
        return nil, t
    }
    p.m.kind = reflect.Int64
    p.m.value = reflect.ValueOf(i)
	return parseEOF, p.next()
}

func parseComplex(p *parser, t Token) (parseFunc, Token) {
    c , err := strconv.ParseComplex(t.Value, 128)
    if err != nil {
        // error handling
        return nil, t
    }
    p.m.kind = reflect.Complex128
    p.m.value = reflect.ValueOf(c)
	return parseEOF, p.next()
}

func parseSlice(p *parser, t Token) (parseFunc, Token) {
	return nil, t
}

func parseEOF(p *parser, _ Token) (parseFunc, Token) {
	p.emit()
	return parsePlus, p.next()
}
