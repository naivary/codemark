package marker

import (
	"reflect"
)

type Marker interface {
	// Ident is the identifier of
	// the marker without the `+`
	Ident() string

	// Kind is the reflect.Kind
	// the marker is using
	Kind() reflect.Kind

	Value() reflect.Value
}

func NewDefault(ident string, kind reflect.Kind, value reflect.Value) *Default {
	return &Default{
		Idn: ident,
		K:   kind,
		Val: value,
	}
}

var _ Marker = (*Default)(nil)

type Default struct {
	Idn string
	K   reflect.Kind
	Val reflect.Value
}

func (d *Default) Ident() string {
	return d.Idn
}

func (d *Default) Kind() reflect.Kind {
	return d.K
}

func (d *Default) Value() reflect.Value {
	return d.Val
}
