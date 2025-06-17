package utils

import (
	"reflect"

	"github.com/naivary/codemark/parser"
)

func MarkerKindOf(typ reflect.Type) parser.MarkerKind {
	kind := Deref(typ).Kind()
	switch kind {
	case reflect.Slice:
		return parser.MarkerKindList
	// rune=uint8 & byte=int32
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return parser.MarkerKindInt
	case reflect.Float32, reflect.Float64:
		return parser.MarkerKindFloat
	case reflect.Complex64, reflect.Complex128:
		return parser.MarkerKindComplex
	case reflect.Bool:
		return parser.MarkerKindBool
	case reflect.String:
		return parser.MarkerKindString
	}
	return 0
}
