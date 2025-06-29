package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/syntax"
)

type Marker struct {
	ident string
	kind  marker.Kind
	value reflect.Value
}

// NewMarker returns a new Marker WITHOUT any validations. If you want to create
// a marker which is getting validated use `maker.MakeMarker`.
func NewMarker(ident string, kind marker.Kind, value reflect.Value) Marker {
	return Marker{
		ident: ident,
		kind:  kind,
		value: value,
	}
}

func (m *Marker) String() string {
	if m.kind == marker.STRING {
		return fmt.Sprintf(`%s="%v"`, m.ident, m.value)
	}
	if m.kind == marker.LIST {
		list := fmt.Sprintf(`%#v`, m.value)
		list, _ = strings.CutPrefix(list, "[]interface {}")
		list = strings.ReplaceAll(list, "{", "[")
		list = strings.ReplaceAll(list, "}", "]")
		list = strings.ReplaceAll(list, "(", "")
		list = strings.ReplaceAll(list, ")", "")
		return fmt.Sprintf("%s=%s", m.ident, list)
	}
	if m.kind == marker.COMPLEX {
		c := fmt.Sprintf("%v", m.value)
		c = strings.ReplaceAll(c, "(", "")
		c = strings.ReplaceAll(c, ")", "")
		return fmt.Sprintf("%s=%v", m.ident, c)
	}
	return fmt.Sprintf("%s=%v", m.ident, m.value)
}

func (m *Marker) Ident() string {
	return m.ident
}

func (m *Marker) Kind() marker.Kind {
	return m.kind
}

func (m *Marker) Value() reflect.Value {
	return m.value
}

func (m *Marker) IsValid() error {
	if err := syntax.Ident(m.ident); err != nil {
		return fmt.Errorf("marker identifier is invalid: %s\n", m.ident)
	}
	if !m.value.IsValid() {
		return fmt.Errorf("value of markeris not valid: %v\n", m.value)
	}
	return nil
}
