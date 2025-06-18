package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/naivary/codemark/parser/marker"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

type Marker struct {
	ident string
	kind  marker.Kind
	value reflect.Value
}

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
	return sdkutil.IsValidIdent(m.ident)
}
