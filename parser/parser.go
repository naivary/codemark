package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/naivary/codemark/lexer"
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

func Parse(input string) ([]Marker, error) {
	const minMarker = 1
	p := &parser{
		l:       lexer.Lex(input),
		state:   parsePlus,
		markers: make([]Marker, 0, minMarker),
		m:       &marker{},
	}
	p.run()
	return p.markers, p.err
}

type parser struct {
	l *lexer.Lexer

	state parseFunc
	// all the markers which have been built, including error markers
	markers []Marker

	// the current marker which is being built
	m *marker

	err error
}

func (p *parser) run() {
	next := false
	token := p.l.NextToken()
	for state := p.state; state != nil; {
		if next {
			token = p.l.NextToken()
		}
		if token.Kind == lexer.TokenKindError {
			state, next = p.errorf("failed while lexing: %s", token.Value)
			// we can either break or use continue. To convey to the usual
			// pattern of returning the next state and _next or _keep it's
			// better to use conutinue
			continue
		}
		state, next = state(p, token)
	}
}

func (p *parser) emit() {
	p.markers = append(p.markers, p.m)
	p.m = &marker{}
}

func (p *parser) errorf(format string, args ...any) (parseFunc, bool) {
	p.err = fmt.Errorf(format, args...)
	return nil, _keep
}

func parsePlus(p *parser, t lexer.Token) (parseFunc, bool) {
	if t.Kind == lexer.TokenKindEOF {
		return nil, _keep
	}
	return parseIdent, _next
}

func parseIdent(p *parser, t lexer.Token) (parseFunc, bool) {
	p.m.ident = t.Value
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
		return parseListStart, _keep
	default:
		return p.errorf("A wrong kind is passed as a TokenKind from the lexer. This should usually never happen! Found kind is: `%s`", t)
	}
}

func parseBool(p *parser, t lexer.Token) (parseFunc, bool) {
	val, err := strconv.ParseBool(t.Value)
	if err != nil {
		return p.errorf("couldn't parse boolean value: %s", t.Value)
	}
	p.m.kind = MarkerKindBool
	p.m.value = reflect.ValueOf(val)
	return parseEOF, _next
}

func parseString(p *parser, t lexer.Token) (parseFunc, bool) {
	p.m.kind = MarkerKindString
	p.m.value = reflect.ValueOf(t.Value)
	return parseEOF, _next
}

func parseInt(p *parser, t lexer.Token) (parseFunc, bool) {
	val, err := parseInt64(t.Value)
	if err != nil {
		return p.errorf("couldn't parse int value: `%s`. Err: %v", t.Value, err)
	}
	p.m.kind = MarkerKindInt
	p.m.value = reflect.ValueOf(val)
	return parseEOF, _next
}

func parseFloat(p *parser, t lexer.Token) (parseFunc, bool) {
	val, err := parseFloat64(t.Value)
	if err != nil {
		return p.errorf("couldn't parse float value: `%s`. Err: %v", t.Value, err)
	}
	p.m.kind = MarkerKindFloat
	p.m.value = reflect.ValueOf(val)
	return parseEOF, _next
}

func parseComplex(p *parser, t lexer.Token) (parseFunc, bool) {
	c, err := parseComplex128(t.Value)
	if err != nil {
		return p.errorf("couldn't parse complex value: `%s`. Err: %v", t.Value, err)
	}
	p.m.kind = MarkerKindComplex
	p.m.value = reflect.ValueOf(c)
	return parseEOF, _next
}

func parseListStart(p *parser, t lexer.Token) (parseFunc, bool) {
	p.m.kind = MarkerKindList
	// slice of type any because the choosen output type of the user might be
	// []any.
	rtype := reflect.TypeOf([]any{})
	p.m.value = reflect.MakeSlice(rtype, 0, 1)
	return parseListElem, _next
}

func parseListElem(p *parser, t lexer.Token) (parseFunc, bool) {
	switch t.Kind {
	case lexer.TokenKindString:
		return parseListStringElem, _keep
	case lexer.TokenKindInt:
		return parseListIntElem, _keep
	case lexer.TokenKindFloat:
		return parseListFloatElem, _keep
	case lexer.TokenKindComplex:
		return parseListComplexElem, _keep
	case lexer.TokenKindBool:
		return parseListBoolElem, _keep
	default:
		return parseListEnd, _keep
	}
}

func parseListBoolElem(p *parser, t lexer.Token) (parseFunc, bool) {
	val := reflect.ValueOf(t.Value)
	p.m.value = reflect.Append(p.m.value, val)
	return parseListElem, _next
}

func parseListStringElem(p *parser, t lexer.Token) (parseFunc, bool) {
	val := reflect.ValueOf(t.Value)
	p.m.value = reflect.Append(p.m.value, val)
	return parseListElem, _next
}

func parseListIntElem(p *parser, t lexer.Token) (parseFunc, bool) {
	i, err := parseInt64(t.Value)
	if err != nil {
		return p.errorf("couldn't parse int value `%s` in array. Err: %v", t.Value, err)
	}
	val := reflect.ValueOf(i)
	p.m.value = reflect.Append(p.m.value, val)
	return parseListElem, _next
}

func parseListFloatElem(p *parser, t lexer.Token) (parseFunc, bool) {
	f, err := parseFloat64(t.Value)
	if err != nil {
		return p.errorf("couldn't parse float value `%s` in array. Err: %v", t.Value, err)
	}
	val := reflect.ValueOf(f)
	p.m.value = reflect.Append(p.m.value, val)
	return parseListElem, _next
}

func parseListComplexElem(p *parser, t lexer.Token) (parseFunc, bool) {
	c, err := parseComplex128(t.Value)
	if err != nil {
		return p.errorf("couldn't parse complex value `%s` in array. Err: %v", t.Value, err)
	}
	val := reflect.ValueOf(c)
	p.m.value = reflect.Append(p.m.value, val)
	return parseListElem, _next
}

func parseListEnd(p *parser, _ lexer.Token) (parseFunc, bool) {
	return parseEOF, _next
}

func parseEOF(p *parser, _ lexer.Token) (parseFunc, bool) {
	p.emit()
	return parsePlus, _keep
}
