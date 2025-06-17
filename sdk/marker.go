package sdk

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/naivary/codemark/lexer"
	"github.com/naivary/codemark/parser"
)

type Marker interface {
	String() string
	// Ident is the identifier of the marker without `+`
	Ident() string

	// Kind of the Marker
	Kind() parser.MarkerKind

	// Value of the marker defined on the right side of the assignment `=`
	Value() reflect.Value
}

var _ Marker = (*marker)(nil)

type marker struct {
	ident string
	kind  parser.MarkerKind
	value reflect.Value
}

func NewMarker(ident string, kind parser.MarkerKind, value reflect.Value) Marker {
	m := &marker{
		ident: ident,
		kind:  kind,
		value: value,
	}
	return m
}

func (m *marker) String() string {
	if m.kind == parser.MarkerKindString {
		return fmt.Sprintf(`%s="%v"`, m.ident, m.value)
	}
	if m.kind == parser.MarkerKindList {
		list := fmt.Sprintf(`%#v`, m.value)
		list, _ = strings.CutPrefix(list, "[]interface {}")
		list = strings.ReplaceAll(list, "{", "[")
		list = strings.ReplaceAll(list, "}", "]")
		list = strings.ReplaceAll(list, "(", "")
		list = strings.ReplaceAll(list, ")", "")
		return fmt.Sprintf("%s=%s", m.ident, list)
	}
	if m.kind == parser.MarkerKindComplex {
		c := fmt.Sprintf("%v", m.value)
		c = strings.ReplaceAll(c, "(", "")
		c = strings.ReplaceAll(c, ")", "")
		return fmt.Sprintf("%s=%v", m.ident, c)
	}
	return fmt.Sprintf("%s=%v", m.ident, m.value)
}

func (m *marker) Ident() string {
	return m.ident
}

func (m *marker) Kind() parser.MarkerKind {
	return m.kind
}

func (m *marker) Value() reflect.Value {
	return m.value
}

func (m *marker) IsValid() error {
	return lexer.IsValidIdent(m.ident)
}
