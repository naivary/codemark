package marker

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/naivary/codemark/internal/equal"
	"github.com/naivary/codemark/syntax"
)

// Fake returns a regular marker with the identifier prefix of
// "codemark:fake:*". This is useful if you want to call another converter
// inside a converter.
func Fake(mkind Kind, value reflect.Value) Marker {
	ident := fmt.Sprintf("codemark:fake:%s", mkind.String())
	return New(ident, mkind, value)
}

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

// String returns the marker as it would be written in a go comment.
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

// IsValid validates the marker. The marker is valid if the returned error is
// nil.
func (m *Marker) IsValid() error {
	if err := syntax.Ident(m.Ident); err != nil {
		return fmt.Errorf("marker identifier is invalid: %s", m.Ident)
	}
	if !m.Value.IsValid() {
		return fmt.Errorf("value of marker is not valid: %v", m.Value)
	}
	return nil
}

// IsEqual checks if the given reflect.Value is equal to the marker value. The
// given value must be of the same kind of the marker e.g. a STRING marker can
// only be compared to a value which is also a string e.g. the method
// `reflect.Value.String()` will not fail. So before providing the value to the
// function make sure it is of the type expected.
func (m *Marker) IsEqual(v reflect.Value) bool {
	return equal.IsEqual(v, m.Value)
}
