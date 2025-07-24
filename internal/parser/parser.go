package parser

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/naivary/codemark/internal/lexer"
	"github.com/naivary/codemark/internal/lexer/token"
	"github.com/naivary/codemark/internal/utils"
	"github.com/naivary/codemark/marker"
)

var isListSeq = false

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
	c, err := strconv.ParseComplex(utils.ComplexOrder(v), 128)
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

func Parse(input string) ([]marker.Marker, error) {
	const minMarker = 1
	p := &parser{
		l:       lexer.Lex(input),
		state:   parsePlus,
		markers: make([]marker.Marker, 0, minMarker),
		m:       &marker.Marker{},
	}
	p.run()
	return p.markers, p.err
}

type parser struct {
	l *lexer.Lexer

	state parseFunc

	// all the markers which have been built, including error markers
	markers []marker.Marker

	// the current marker which is being built
	m *marker.Marker

	err error
}

func (p *parser) run() {
	next := false
	t := p.l.NextToken()
	for state := p.state; state != nil; {
		if next {
			t = p.l.NextToken()
		}
		if t.Kind == token.ERROR {
			state, next = p.errorf("failed while lexing: %s", t.Value)
			// we can either break or use continue. To convey to the usual
			// pattern of returning the next state and _next or _keep it's
			// better to use continue
			continue
		}
		state, next = state(p, t)
	}
}

func (p *parser) emit() {
	p.markers = append(p.markers, *p.m)
	p.m = &marker.Marker{}
}

func (p *parser) errorf(format string, args ...any) (parseFunc, bool) {
	p.err = fmt.Errorf(format, args...)
	return nil, _keep
}

func parsePlus(p *parser, t lexer.Token) (parseFunc, bool) {
	if t.Kind == token.EOF {
		return nil, _keep
	}
	return parseIdent, _next
}

func parseIdent(p *parser, t lexer.Token) (parseFunc, bool) {
	p.m.Ident = t.Value
	return parseAssign, _next
}

func parseAssign(p *parser, t lexer.Token) (parseFunc, bool) {
	return parseValue, _next
}

func parseValue(p *parser, t lexer.Token) (parseFunc, bool) {
	switch t.Kind {
	case token.STRING:
		return parseString, _keep
	case token.INT:
		return parseInt, _keep
	case token.COMPLEX:
		return parseComplex, _keep
	case token.FLOAT:
		return parseFloat, _keep
	case token.BOOL:
		return parseBool, _keep
	case token.LBRACK:
		return parseList, _keep
	default:
		return p.errorf(
			"A wrong kind is passed as a TokenKind from the lexer. This should usually never happen! Found kind is: `%s`",
			t,
		)
	}
}

func parseBool(p *parser, t lexer.Token) (parseFunc, bool) {
	val, err := strconv.ParseBool(t.Value)
	if err != nil {
		return p.errorf("couldn't parse boolean value: %s", t.Value)
	}
	rvalue := reflect.ValueOf(val)
	if isListSeq {
		p.m.Value = reflect.Append(p.m.Value, rvalue)
		return parseListElem, _next
	}
	p.m.Kind = marker.BOOL
	p.m.Value = rvalue
	return parseEOF, _next
}

func parseString(p *parser, t lexer.Token) (parseFunc, bool) {
	rvalue := reflect.ValueOf(t.Value)
	if isListSeq {
		p.m.Value = reflect.Append(p.m.Value, rvalue)
		return parseListElem, _next
	}

	p.m.Kind = marker.STRING
	p.m.Value = rvalue
	return parseEOF, _next
}

func parseInt(p *parser, t lexer.Token) (parseFunc, bool) {
	val, err := parseInt64(t.Value)
	if err != nil {
		return p.errorf("couldn't parse int value: `%s`. Err: %v", t.Value, err)
	}
	rvalue := reflect.ValueOf(val)
	if isListSeq {
		p.m.Value = reflect.Append(p.m.Value, rvalue)
		return parseListElem, _next
	}
	p.m.Kind = marker.INT
	p.m.Value = rvalue
	return parseEOF, _next
}

func parseFloat(p *parser, t lexer.Token) (parseFunc, bool) {
	val, err := parseFloat64(t.Value)
	if err != nil {
		return p.errorf("couldn't parse float value: `%s`. Err: %v", t.Value, err)
	}
	rvalue := reflect.ValueOf(val)
	if isListSeq {
		p.m.Value = reflect.Append(p.m.Value, rvalue)
		return parseListElem, _next
	}
	p.m.Kind = marker.FLOAT
	p.m.Value = reflect.ValueOf(val)
	return parseEOF, _next
}

func parseComplex(p *parser, t lexer.Token) (parseFunc, bool) {
	val, err := parseComplex128(t.Value)
	if err != nil {
		return p.errorf("couldn't parse complex value: `%s`. Err: %v", t.Value, err)
	}
	rvalue := reflect.ValueOf(val)
	if isListSeq {
		p.m.Value = reflect.Append(p.m.Value, rvalue)
		return parseListElem, _next
	}
	p.m.Kind = marker.COMPLEX
	p.m.Value = rvalue
	return parseEOF, _next
}

func parseList(p *parser, t lexer.Token) (parseFunc, bool) {
	p.m.Kind = marker.LIST
	// slice must be of type any because the choosen output type of the
	// user might be []any.
	rtype := reflect.TypeOf([]any{})
	p.m.Value = reflect.MakeSlice(rtype, 0, 1)
	isListSeq = true
	return parseListElem, _next
}

func parseListElem(p *parser, t lexer.Token) (parseFunc, bool) {
	switch t.Kind {
	case token.STRING:
		return parseString, _keep
	case token.INT:
		return parseInt, _keep
	case token.FLOAT:
		return parseFloat, _keep
	case token.COMPLEX:
		return parseComplex, _keep
	case token.BOOL:
		return parseBool, _keep
	default:
		return parseListEnd, _keep
	}
}

func parseListEnd(p *parser, _ lexer.Token) (parseFunc, bool) {
	isListSeq = false
	return parseEOF, _next
}

func parseEOF(p *parser, _ lexer.Token) (parseFunc, bool) {
	p.emit()
	return parsePlus, _keep
}
