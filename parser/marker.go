package parser

import (
	"reflect"
)

type MarkerKind int

const (
	MarkerKindList = iota + 1
	MarkerKindFloat
	MarkerKindInt
	MarkerKindComplex
	MarkerKindBool
	MarkerKindString
)

type Marker interface {
	// Ident is the identifier of the marker without `+`
	Ident() string

	// Kind of the Marker
	Kind() MarkerKind

	// Value of the marker defined on the right side of the assignment `=`
	Value() reflect.Value
}

func NewMarker(ident string, kind MarkerKind, value reflect.Value) Marker {
	return &marker{
		Idn: ident,
		K:   kind,
		Val: value,
	}
}

var _ Marker = (*marker)(nil)

type marker struct {
	Idn string
	K   MarkerKind
	Val reflect.Value
}

func (m *marker) Ident() string {
	return m.Idn
}

func (m *marker) Kind() MarkerKind {
	return m.K
}

func (m *marker) Value() reflect.Value {
	return m.Val
}
