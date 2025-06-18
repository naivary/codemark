package sdk

import (
	"reflect"
)

type Marker interface {
	String() string
	// Ident is the identifier of the marker without `+`
	Ident() string

	// Kind of the Marker
	Kind() MarkerKind

	// Value of the marker defined on the right side of the assignment `=`
	Value() reflect.Value
}
