package parser

import (
	"fmt"
	"reflect"
)

func MarkerKindOf(typ reflect.Type) MarkerKind {
	kind := typ.Kind()
	if kind == reflect.Ptr {
		kind = typ.Elem().Kind()
	}
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

var markerNames = map[MarkerKind]string{
	MarkerKindString:  "MarkerKindString",
	MarkerKindFloat:   "MarkerKindFloat",
	MarkerKindInt:     "MarkerKindInt",
	MarkerKindComplex: "MarkerKindComplex",
	MarkerKindBool:    "MarkerKindBool",
	MarkerKindList:    "MarkerKindList",
}

func (m MarkerKind) String() string {
	if name, ok := markerNames[m]; ok {
		return name
	}
	return fmt.Sprintf("MarkerKind<%d>", m)
}

type Marker interface {
	String() string
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

func (m *marker) String() string {
	return fmt.Sprintf("%s:%v", m.Idn, m.Val)
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
