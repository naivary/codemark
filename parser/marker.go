package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Marker = (*marker)(nil)

type marker struct {
	ident string
	kind  sdk.MarkerKind
	value reflect.Value
}

func NewMarker(ident string, kind sdk.MarkerKind, value reflect.Value) sdk.Marker {
	m := &marker{
		ident: ident,
		kind:  kind,
		value: value,
	}
	return m
}

func (m *marker) String() string {
	if m.kind == sdk.MarkerKindString {
		return fmt.Sprintf(`%s="%v"`, m.ident, m.value)
	}
	if m.kind == sdk.MarkerKindList {
		list := fmt.Sprintf(`%#v`, m.value)
		list, _ = strings.CutPrefix(list, "[]interface {}")
		list = strings.ReplaceAll(list, "{", "[")
		list = strings.ReplaceAll(list, "}", "]")
		list = strings.ReplaceAll(list, "(", "")
		list = strings.ReplaceAll(list, ")", "")
		return fmt.Sprintf("%s=%s", m.ident, list)
	}
	if m.kind == sdk.MarkerKindComplex {
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

func (m *marker) Kind() sdk.MarkerKind {
	return m.kind
}

func (m *marker) Value() reflect.Value {
	return m.value
}

func (m *marker) IsValid() error {
	return sdkutil.IsValidIdent(m.ident)
}
