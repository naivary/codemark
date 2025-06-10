package codemark

import (
	"reflect"
	"strings"
)

const (
	_rune = reflect.Int32
	_byte = reflect.Uint8

	TypeIDSep = "."
)

var (
	_rvzero = reflect.Value{}
)

func deref(typ reflect.Type) reflect.Type {
	if isPointer(typ) {
		typ = typ.Elem()
	}
	return typ
}

func isPointer(typ reflect.Type) bool {
	return typ.Kind() == reflect.Pointer
}

func toOutput(v reflect.Value, typ reflect.Type) (reflect.Value, error) {
	isPtr := isPointer(typ)
	// have to create a new variable because the original type might be needed
	// in case of a pointer.
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

func typeID(typ reflect.Type, b *strings.Builder) string {
	if typ == nil {
		return b.String()
	}
	// write dot to seperate the incoming kind
	if len(b.String()) > 0 {
		b.WriteString(TypeIDSep)
	}
	kind := typ.Kind()
	if _, err := b.WriteString(kind.String()); err != nil {
		panic(err)
	}
	if kind == reflect.Pointer {
		return typeID(typ.Elem(), b)
	}
	if kind == reflect.Map {
		var keyBuilder strings.Builder
		key := typeID(typ.Key(), &keyBuilder)
		b.WriteString(TypeIDSep)
		b.WriteString(key)
		return typeID(typ.Elem(), b)
	}
	if kind == reflect.Slice {
		return typeID(typ.Elem(), b)
	}
	if kind == reflect.Chan {
		return typeID(typ.Elem(), b)
	}
	if kind == reflect.Array {
		return typeID(typ.Elem(), b)
	}
	return typeID(nil, b)
}

func TypeID(typ reflect.Type) string {
	if typ == nil {
		return ""
	}
	var b strings.Builder
	return typeID(typ, &b)
}

func TypeIDFromAny(typ any) string {
	if typ == nil {
		return ""
	}
	var b strings.Builder
	rtype := reflect.TypeOf(typ)
	return typeID(rtype, &b)
}
