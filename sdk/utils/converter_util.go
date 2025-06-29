package utils

import (
	"fmt"
	"reflect"
	"strings"
)

const NameSep = "."

// NewConvName is returning a valid converter name. The convention is to prefix
// every converter with your project name, followed by a custom name for the
// converter seperated by a dot.
func NewConvName(proj, conv string) string {
	return fmt.Sprintf("%s.%s", proj, conv)
}

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

// ConvertTo converts the value `v` to the Type `typ` handling pointer
// dereferencing and other inconveniences.
func ConvertTo(v reflect.Value, typ reflect.Type) (reflect.Value, error) {
	isPtr := IsPointer(typ)
	outputType := typ
	// need to dereference type to create the correct variable using
	// `reflect.New`. Otherwise .Set wont work.
	if isPtr {
		outputType = outputType.Elem()
	}
	out := reflect.New(outputType)
	out.Elem().Set(v.Convert(outputType))
	if isPtr {
		return out.Convert(typ), nil
	}
	return out.Elem(), nil
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
