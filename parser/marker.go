package parser

import "reflect"

type Marker interface {
	// Ident is the identifier of
	// the marker without the `+`
	Ident() string

	// Kind is the reflect.Kind
	// the marker is using
	Kind() reflect.Kind

	Value() reflect.Value
}

func NewMarker(ident string, kind reflect.Kind, value reflect.Value) Marker {
	return &marker{
		ident: ident,
		kind:  kind,
		value: value,
	}
}

var _ Marker = (*marker)(nil)

type marker struct {
	ident string
	kind  reflect.Kind
	value reflect.Value
}

func (m *marker) Ident() string {
	return m.ident
}

func (m *marker) Kind() reflect.Kind {
	return m.kind
}

func (m *marker) Value() reflect.Value {
	return m.value
}
