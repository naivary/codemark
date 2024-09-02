package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/naivary/codemark/lexer"
	"github.com/naivary/codemark/marker"
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
type parseFunc func(*parser, lexer.Token) (parseFunc, bool)

func parse(input string) *parser {
	return &parser{
		l:       lexer.Lex(input),
		state:   parsePlus,
		markers: make(chan marker.Marker, 2),
		m:       &marker.Default{},
	}
}

type parser struct {
	l       *lexer.Lexer
	state   parseFunc
	markers chan marker.Marker
	m       *marker.Default
}

func (p *parser) run() {
	p.l.Run()
	next := false
	token := p.l.NextToken()
	for state := p.state; state != nil; {
		if next {
			token = p.l.NextToken()
		}
		if token.Kind == lexer.TokenKindError {
			state, next = p.errorf("failed while lexing: %s", token.Value)
			return
		}
		state, next = state(p, token)
	}
	close(p.markers)
}

func (p *parser) emit() {
	p.markers <- p.m
	p.m = &marker.Default{}
}

func (p *parser) errorf(format string, args ...any) (parseFunc, bool) {
	err := fmt.Errorf(format, args...)
	p.m.SetError(err)
	p.emit()
	close(p.markers)
	return nil, _keep
}

func parsePlus(p *parser, t lexer.Token) (parseFunc, bool) {
	if t.Kind == lexer.TokenKindEOF {
		return nil, _keep
	}
	return parseIdent, _next
}

func parseIdent(p *parser, t lexer.Token) (parseFunc, bool) {
	p.m.Idn = t.Value
	return parseAssignment, _next
}

func parseAssignment(p *parser, t lexer.Token) (parseFunc, bool) {
	return parseValue, _next
}

func parseValue(p *parser, t lexer.Token) (parseFunc, bool) {
	switch t.Kind {
	case lexer.TokenKindString:
		return parseString, _keep
	case lexer.TokenKindInt:
		return parseInt, _keep
	case lexer.TokenKindComplex:
		return parseComplex, _keep
	case lexer.TokenKindFloat:
		return parseFloat, _keep
	case lexer.TokenKindBool:
		return parseBool, _keep
	case lexer.TokenKindOpenSquareBracket:
		return parseSliceStart, _keep
	default:
		return p.errorf("A wrong kind is passed as a TokenKind from the lexer. This should usually never happen! Found kind is: `%s`", t)
	}
}

func parseBool(p *parser, t lexer.Token) (parseFunc, bool) {
	p.m.K = reflect.Bool
	boolVal, err := strconv.ParseBool(t.Value)
	if err != nil {
		return p.errorf("couldn't parse boolean value: %s", t.Value)
	}
	p.m.Val = reflect.ValueOf(boolVal)
	return parseEOF, _next
}

func parseString(p *parser, t lexer.Token) (parseFunc, bool) {
	p.m.K = reflect.String
	p.m.Val = reflect.ValueOf(t.Value)
	return parseEOF, _next
}

func parseInt(p *parser, t lexer.Token) (parseFunc, bool) {
	i, err := parseInt64(t.Value)
	if err != nil {
		return p.errorf("couldn't parse int value: `%s`. Err: %v", t.Value, err)
	}
	p.m.K = reflect.Int64
	p.m.Val = reflect.ValueOf(i)
	return parseEOF, _next
}

func parseFloat(p *parser, t lexer.Token) (parseFunc, bool) {
	f, err := parseFloat64(t.Value)
	if err != nil {
		return p.errorf("couldn't parse float value: `%s`. Err: %v", t.Value, err)
	}
	p.m.K = reflect.Float64
	p.m.Val = reflect.ValueOf(f)
	return parseEOF, _next
}

func parseComplex(p *parser, t lexer.Token) (parseFunc, bool) {
	c, err := parseComplex128(t.Value)
	if err != nil {
		return p.errorf("couldn't parse complex value: `%s`. Err: %v", t.Value, err)
	}
	p.m.K = reflect.Complex128
	p.m.Val = reflect.ValueOf(c)
	return parseEOF, _next
}

func parseSliceStart(p *parser, t lexer.Token) (parseFunc, bool) {
	p.m.K = reflect.Slice
	rt := reflect.TypeOf([]any{})
	p.m.Val = reflect.MakeSlice(rt, 0, 1)
	return parseSliceElem, _next
}

func parseSliceElem(p *parser, t lexer.Token) (parseFunc, bool) {
	switch t.Kind {
	case lexer.TokenKindString:
		return parseSliceStringElem, _keep
	case lexer.TokenKindInt:
		return parseSliceIntElem, _keep
	case lexer.TokenKindFloat:
		return parseSliceFloatElem, _keep
	case lexer.TokenKindComplex:
		return parseSliceComplexElem, _keep
	case lexer.TokenKindBool:
		return parseSliceBoolElem, _keep
	default:
		return parseSliceEnd, _keep
	}
}

func parseSliceBoolElem(p *parser, t lexer.Token) (parseFunc, bool) {
	val := reflect.ValueOf(t.Value)
	p.m.Val = reflect.Append(p.m.Val, val)
	return parseSliceElem, _next
}

func parseSliceStringElem(p *parser, t lexer.Token) (parseFunc, bool) {
	val := reflect.ValueOf(t.Value)
	p.m.Val = reflect.Append(p.m.Val, val)
	return parseSliceElem, _next
}

func parseSliceIntElem(p *parser, t lexer.Token) (parseFunc, bool) {
	i, err := parseInt64(t.Value)
	if err != nil {
		return p.errorf("couldn't parse int value `%s` in array. Err: %v", t.Value, err)
	}
	val := reflect.ValueOf(i)
	p.m.Val = reflect.Append(p.m.Val, val)
	return parseSliceElem, _next
}

func parseSliceFloatElem(p *parser, t lexer.Token) (parseFunc, bool) {
	f, err := parseFloat64(t.Value)
	if err != nil {
		return p.errorf("couldn't parse float value `%s` in array. Err: %v", t.Value, err)
	}
	val := reflect.ValueOf(f)
	p.m.Val = reflect.Append(p.m.Val, val)
	return parseSliceElem, _next
}

func parseSliceComplexElem(p *parser, t lexer.Token) (parseFunc, bool) {
	c, err := parseComplex128(t.Value)
	if err != nil {
		return p.errorf("couldn't parse complex value `%s` in array. Err: %v", t.Value, err)
	}
	val := reflect.ValueOf(c)
	p.m.Val = reflect.Append(p.m.Val, val)
	return parseSliceElem, _next
}

func parseSliceEnd(p *parser, _ lexer.Token) (parseFunc, bool) {
	return parseEOF, _next
}

func parseEOF(p *parser, _ lexer.Token) (parseFunc, bool) {
	p.emit()
	return parsePlus, _keep
}
