package marker

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/naivary/codemark/syntax"
)

type Marker struct {
	Ident string
	Kind  Kind
	Value reflect.Value
}

// NewMarker returns a new Marker WITHOUT any validations. If you want to create
// a marker which is getting validated use `maker.MakeMarker`.
func New(ident string, kind Kind, value reflect.Value) Marker {
	return Marker{
		Ident: ident,
		Kind:  kind,
		Value: value,
	}
}

func (m *Marker) String() string {
	if m.Kind == STRING {
		return fmt.Sprintf(`%s="%v"`, m.Ident, m.Value)
	}
	if m.Kind == LIST {
		list := fmt.Sprintf(`%#v`, m.Value)
		list, _ = strings.CutPrefix(list, "[]interface {}")
		list = strings.ReplaceAll(list, "{", "[")
		list = strings.ReplaceAll(list, "}", "]")
		list = strings.ReplaceAll(list, "(", "")
		list = strings.ReplaceAll(list, ")", "")
		return fmt.Sprintf("%s=%s", m.Ident, list)
	}
	if m.Kind == COMPLEX {
		c := fmt.Sprintf("%v", m.Value)
		c = strings.ReplaceAll(c, "(", "")
		c = strings.ReplaceAll(c, ")", "")
		return fmt.Sprintf("%s=%v", m.Ident, c)
	}
	return fmt.Sprintf("%s=%v", m.Ident, m.Value)
}

func (m *Marker) IsValid() error {
	if err := syntax.Ident(m.Ident); err != nil {
		return fmt.Errorf("marker identifier is invalid: %s\n", m.Ident)
	}
	if !m.Value.IsValid() {
		return fmt.Errorf("value of markeris not valid: %v\n", m.Value)
	}
	return nil
}
