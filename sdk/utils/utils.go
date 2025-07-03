package utils

import (
	"reflect"
	"slices"

	"github.com/naivary/codemark/parser/marker"
)

// TypeForMarkerKind returns the reflect.Type used for the given marker kind
func TypeForMarkerKind(mkind marker.Kind) reflect.Type {
	switch mkind {
	case marker.STRING:
		return reflect.TypeFor[string]()
	case marker.INT:
		return reflect.TypeFor[int64]()
	case marker.FLOAT:
		return reflect.TypeFor[float64]()
	case marker.COMPLEX:
		return reflect.TypeFor[complex128]()
	case marker.BOOL:
		return reflect.TypeFor[bool]()
	case marker.LIST:
		reflect.TypeFor[[]any]()
	}
	return nil
}

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

func IsAny(rtype reflect.Type) bool {
	return Deref(rtype).Kind() == reflect.Interface
}

// IsPrimitive is returning true iff the given type is non-slice and a type
// which can be converted by a builtin converter.
func IsPrimitive(rtype reflect.Type) bool {
	return IsInt(rtype) || IsUint(rtype) || IsFloat(rtype) || IsString(rtype) || IsBool(rtype) || IsComplex(rtype)
}

func IsValidSlice(rtype reflect.Type) bool {
	if rtype.Kind() != reflect.Slice {
		return false
	}
	elem := rtype.Elem()
	return IsPrimitive(elem) || IsAny(elem)
}
