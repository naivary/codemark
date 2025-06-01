package marker

import (
	"reflect"
)

type Marker interface {
	// Ident is the identifier of the marker without `+`
	Ident() string

	// Kind of the Marker
	Kind() MarkerKind

	// Value of the marker defined on the right side of the assignment `=`
	Value() reflect.Value
}

func NewDefault(ident string, kind MarkerKind, value reflect.Value) *Default {
	return &Default{
		Idn: ident,
		K:   kind,
		Val: value,
	}
}

var _ Marker = (*Default)(nil)

type Default struct {
	Idn string
	K   MarkerKind
	Val reflect.Value
}

func (d *Default) Ident() string {
	return d.Idn
}

func (d *Default) Kind() MarkerKind {
	return d.K
}

func (d *Default) Value() reflect.Value {
	return d.Val
}
