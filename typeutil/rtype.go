package typeutil

import (
	"reflect"
	"slices"
	"strings"
)

const NameSep = "."

// Deref dereferences the type. If the given type is not a pointer it will be
// returned without any dereferencing.
func Deref(typ reflect.Type) reflect.Type {
	if typ == nil {
		return typ
	}
	if IsPointer(typ) {
		typ = typ.Elem()
	}
	return typ
}

func IsPointer(typ reflect.Type) bool {
	return typ.Kind() == reflect.Pointer
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

// NameFor returns a string representation of a given reflect.Type. It is NOT
// unique. The returned string is like reading the type from left to right. Some
// examples are:
// map[string]string -> map.string.string
// []int -> slice.int
// *int -> ptr.int
// The uniqueness is not given because some types are aliases of others e.g:
// *byte -> ptr.uint8
// *uint8 -> ptr.uint8
// and others can not be uniquely identified like structs and interfaces.
func NameFor(rtype reflect.Type) string {
	var b strings.Builder
	return nameFor(rtype, &b)
}

func nameFor(typ reflect.Type, b *strings.Builder) string {
	if typ == nil {
		return b.String()
	}
	if len(b.String()) > 0 {
		b.WriteString(NameSep)
	}
	kind := typ.Kind()
	if _, err := b.WriteString(kind.String()); err != nil {
		panic(err)
	}
	if kind == reflect.Pointer {
		return nameFor(typ.Elem(), b)
	}
	if kind == reflect.Map {
		var keyBuilder strings.Builder
		key := nameFor(typ.Key(), &keyBuilder)
		b.WriteString(NameSep)
		b.WriteString(key)
		return nameFor(typ.Elem(), b)
	}
	if kind == reflect.Slice {
		return nameFor(typ.Elem(), b)
	}
	if kind == reflect.Chan {
		return nameFor(typ.Elem(), b)
	}
	if kind == reflect.Array {
		return nameFor(typ.Elem(), b)
	}
	return nameFor(nil, b)
}
