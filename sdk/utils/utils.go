package utils

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unicode"

	"github.com/naivary/codemark/parser/marker"
)

func MarkerKindOf(typ reflect.Type) marker.Kind {
	kind := Deref(typ).Kind()
	switch kind {
	case reflect.Slice:
		return marker.LIST
	// rune=int32 & byte=uint8
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return marker.INT
	case reflect.Float32, reflect.Float64:
		return marker.FLOAT
	case reflect.Complex64, reflect.Complex128:
		return marker.COMPLEX
	case reflect.Bool:
		return marker.BOOL
	case reflect.String:
		return marker.STRING
	}
	return 0
}

func IsValidIdent(ident string) error {
	plus := "+"
	if strings.HasPrefix(ident, plus) {
		return fmt.Errorf("an identifier should not start with a plus. the plus is just like the `var` keyword in a regular programming language: %s\n", ident)
	}
	colon := ":"
	numOfColons := strings.Count(ident, colon)
	if numOfColons < 2 {
		return fmt.Errorf("expected two colons in `%s` but got %d", ident, numOfColons)
	}
	for pathSegment := range strings.SplitSeq(ident, colon) {
		lastChar := rune(pathSegment[len(pathSegment)-1])
		if !(unicode.IsLetter(lastChar) || unicode.IsDigit(lastChar)) {
			return fmt.Errorf("identifier cannot end with an underscore `_` or dot `.`: %s in %s\n", pathSegment, ident)
		}
	}
	return nil
}

func IsBool(rtype reflect.Type) bool {
	return Deref(rtype).Kind() == reflect.Bool
}

func IsString(rtype reflect.Type) bool {
	return Deref(rtype).Kind() == reflect.String
}

func IsInt(rtype reflect.Type) bool {
	kind := Deref(rtype).Kind()
	ints := []reflect.Kind{
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
	}
	return slices.Contains(ints, kind)
}

func IsUint(rtype reflect.Type) bool {
	kind := Deref(rtype).Kind()
	uints := []reflect.Kind{
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
	}
	return slices.Contains(uints, kind)
}

func IsFloat(rtype reflect.Type) bool {
	kind := Deref(rtype).Kind()
	floats := []reflect.Kind{
		reflect.Float32,
		reflect.Float64,
	}
	return slices.Contains(floats, kind)
}

func IsComplex(rtype reflect.Type) bool {
	kind := Deref(rtype).Kind()
	complexes := []reflect.Kind{
		reflect.Complex64,
		reflect.Complex128,
	}
	return slices.Contains(complexes, kind)
}

// IsSupported is returning true iff the given rtype is supported by the default
// converters.
func IsSupported(rtype reflect.Type) bool {
	return IsPrimitive(rtype) || rtype.Kind() == reflect.Slice
}

// IsPrimitive is returning true iff the given type is non-slice and a type
// which can be converted by a builtin converter.
func IsPrimitive(rtype reflect.Type) bool {
	return IsInt(rtype) || IsUint(rtype) || IsFloat(rtype) || IsString(rtype) || IsBool(rtype) || IsComplex(rtype)
}

func IsValidSlice(rtype reflect.Type) bool {
	return rtype.Kind() == reflect.Slice && IsPrimitive(rtype.Elem())
}
