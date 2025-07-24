package typeutil

import (
	"reflect"
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

func DerefValue(v reflect.Value) reflect.Value {
	if IsPointer(v.Type()) {
		return v.Elem()
	}
	return v
}
