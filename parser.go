package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	_next = true
	_keep = false
)

func parseFloat64(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func parseInt64(val string) (int64, error) {
	return strconv.ParseInt(val, 0, 64)
}

func parseComplex128(v string) (complex128, error) {
	re, img, _ := strings.Cut(v, "+")
	i, err := strconv.ParseInt(re, 0, 64)
	if err != nil {
		return 0, err
	}
	val := fmt.Sprintf("%d+%s", i, img)
	c, err := strconv.ParseComplex(val, 128)
	if err != nil {
		return 0, err
	}
	return c, nil
}

// parseFunc is the function signature for parsing.
// the bool return is deciding whether the next token
// from the channel should be received or the last one
// should be passed instead
type parseFunc func(*parser, Token) (parseFunc, bool)

type Marker interface {
	// Ident is the identifier of
	// the marker without the `+`
	Ident() string

	// Kind is the reflect.Kind
	// the marker is using
	Kind() reflect.Kind

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

func (p *parser) run() error {
	p.l.run()
	next := false
	token := p.next()
	for state := p.state; state != nil; {
		if next {
			token = p.next()
		}
		if token.Kind == TokenKindError {
			return errors.New(token.Value)
		}
		state, next = state(p, token)
	}
	close(p.markers)
	return nil
}

func (p *parser) emit() {
	p.markers <- p.m
	p.m = &marker{}
}

func (p *parser) next() Token {
	return <-p.l.tokens
}

func parsePlus(p *parser, t Token) (parseFunc, bool) {
	if t.Kind == TokenKindEOF {
		return nil, _keep
	}
	return parseIdent, _next
}

func parseIdent(p *parser, t Token) (parseFunc, bool) {
	p.m.ident = t.Value
	return parseAssignment, _next
}

func parseAssignment(p *parser, t Token) (parseFunc, bool) {
	return parseValue, _next
}

func parseValue(p *parser, t Token) (parseFunc, bool) {
	switch t.Kind {
	case TokenKindString:
		return parseString, _keep
	case TokenKindInt:
		return parseInt, _keep
	case TokenKindComplex:
		return parseComplex, _keep
	case TokenKindFloat:
		return parseFloat, _keep
	case TokenKindBool:
		return parseBool, _keep
	case TokenKindOpenSquareBracket:
		return parseSliceStart, _keep
	default:
		// Something went wrong while lexing
		return nil, _keep
	}
}

func parseBool(p *parser, t Token) (parseFunc, bool) {
	p.m.kind = reflect.Bool
	boolVal, err := strconv.ParseBool(t.Value)
	if err != nil {
		// error handling
		return nil, _keep
	}
	p.m.value = reflect.ValueOf(boolVal)
	return parseEOF, _next
}

func parseString(p *parser, t Token) (parseFunc, bool) {
	p.m.kind = reflect.String
	p.m.value = reflect.ValueOf(t.Value)
	return parseEOF, _next
}

func parseInt(p *parser, t Token) (parseFunc, bool) {
	i, err := parseInt64(t.Value)
	if err != nil {
		// error handling
		return nil, _keep
	}
	p.m.kind = reflect.Int64
	p.m.value = reflect.ValueOf(i)
	return parseEOF, _next
}

func parseFloat(p *parser, t Token) (parseFunc, bool) {
	f, err := parseFloat64(t.Value)
	if err != nil {
		return nil, _keep
	}
	p.m.kind = reflect.Float64
	p.m.value = reflect.ValueOf(f)
	return parseEOF, _next
}

func parseComplex(p *parser, t Token) (parseFunc, bool) {
	c, err := parseComplex128(t.Value)
	if err != nil {
		return nil, _keep
	}
	p.m.kind = reflect.Complex128
	p.m.value = reflect.ValueOf(c)
	return parseEOF, _next
}

func parseSliceStart(p *parser, t Token) (parseFunc, bool) {
	p.m.kind = reflect.Slice
	rt := reflect.TypeOf([]any{})
	p.m.value = reflect.MakeSlice(rt, 0, 1)
	return parseSliceElem, _next
}

func parseSliceElem(p *parser, t Token) (parseFunc, bool) {
	switch t.Kind {
	case TokenKindString:
		return parseSliceStringElem, _keep
	case TokenKindInt:
		return parseSliceIntElem, _keep
	case TokenKindFloat:
		return parseSliceFloatElem, _keep
	case TokenKindComplex:
		return parseSliceComplexElem, _keep
	default:
		return parseSliceEnd, _keep
	}
}

func parseSliceStringElem(p *parser, t Token) (parseFunc, bool) {
	val := reflect.ValueOf(t.Value)
	p.m.value = reflect.Append(p.m.value, val)
	return parseSliceElem, _next
}

func parseSliceIntElem(p *parser, t Token) (parseFunc, bool) {
	i, err := parseInt64(t.Value)
	if err != nil {
		return nil, _keep
	}
	val := reflect.ValueOf(i)
	p.m.value = reflect.Append(p.m.value, val)
	return parseSliceElem, _next
}

func parseSliceFloatElem(p *parser, t Token) (parseFunc, bool) {
	f, err := parseFloat64(t.Value)
	if err != nil {
		return nil, _keep
	}
	val := reflect.ValueOf(f)
	p.m.value = reflect.Append(p.m.value, val)
	return parseSliceElem, _next
}

func parseSliceComplexElem(p *parser, t Token) (parseFunc, bool) {
	c, err := parseComplex128(t.Value)
	if err != nil {
		return nil, _keep
	}
	val := reflect.ValueOf(c)
	p.m.value = reflect.Append(p.m.value, val)
	return parseSliceElem, _next
}

func parseSliceEnd(p *parser, _ Token) (parseFunc, bool) {
	return parseEOF, _next
}

func parseEOF(p *parser, _ Token) (parseFunc, bool) {
	p.emit()
	return parsePlus, _keep
}
