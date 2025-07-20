package marker

import (
	"reflect"

	"github.com/naivary/codemark/typeutil"
)

// TypeOf returns the reflect.Type used for the given marker kind
func TypeOf(mkind Kind) reflect.Type {
	switch mkind {
	case STRING:
		return reflect.TypeFor[string]()
	case INT:
		return reflect.TypeFor[int64]()
	case FLOAT:
		return reflect.TypeFor[float64]()
	case COMPLEX:
		return reflect.TypeFor[complex128]()
	case BOOL:
		return reflect.TypeFor[bool]()
	case LIST:
		reflect.TypeFor[[]any]()
	}
	return nil
}

// KindOf returns which kind of marker the given type is. If not kind can be
// found 0 will be returned.
func KindOf(typ reflect.Type) Kind {
	switch {
	case typeutil.IsValidSlice(typ):
		return LIST
	case typeutil.IsInt(typ):
		return INT
	case typeutil.IsFloat(typ):
		return FLOAT
	case typeutil.IsComplex(typ):
		return COMPLEX
	case typeutil.IsBool(typ):
		return BOOL
	case typeutil.IsString(typ):
		return STRING
	}
	return 0
}
