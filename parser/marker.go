//go:generate stringer -type=MarkerKind

// TODO: bring marker to sdk and sdk/utils
package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/naivary/codemark/lexer"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

func MarkerKindOf(typ reflect.Type) MarkerKind {
	kind := sdkutil.Deref(typ).Kind()
	switch kind {
	case reflect.Slice:
		return MarkerKindList
	// rune=uint8 & byte=int32
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return MarkerKindInt
	case reflect.Float32, reflect.Float64:
		return MarkerKindFloat
	case reflect.Complex64, reflect.Complex128:
		return MarkerKindComplex
	case reflect.Bool:
		return MarkerKindBool
	case reflect.String:
		return MarkerKindString
	}
	return 0
}

type MarkerKind int

const (
	MarkerKindString MarkerKind = iota + 1
	MarkerKindFloat
	MarkerKindInt
	MarkerKindComplex
	MarkerKindBool
	MarkerKindList
)

type Marker interface {
	String() string
	// Ident is the identifier of the marker without `+`
	Ident() string

	// Kind of the Marker
	Kind() MarkerKind

	// Value of the marker defined on the right side of the assignment `=`
	Value() reflect.Value
}

var _ Marker = (*marker)(nil)

type marker struct {
	ident string
	kind  MarkerKind
	value reflect.Value
}

func NewMarker(ident string, kind MarkerKind, value reflect.Value) Marker {
	m := &marker{
		ident: ident,
		kind:  kind,
		value: value,
	}
	return m
}

func (m *marker) String() string {
	if m.kind == MarkerKindString {
		return fmt.Sprintf(`%s="%v"`, m.ident, m.value)
	}
	if m.kind == MarkerKindList {
		list := fmt.Sprintf(`%#v`, m.value)
		list, _ = strings.CutPrefix(list, "[]interface {}")
		list = strings.ReplaceAll(list, "{", "[")
		list = strings.ReplaceAll(list, "}", "]")
		list = strings.ReplaceAll(list, "(", "")
		list = strings.ReplaceAll(list, ")", "")
		return fmt.Sprintf("%s=%s", m.ident, list)
	}
	if m.kind == MarkerKindComplex {
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

func (m *marker) Kind() MarkerKind {
	return m.kind
}

func (m *marker) Value() reflect.Value {
	return m.value
}

func (m *marker) IsValid() error {
	return lexer.IsValidIdent(m.ident)
}
