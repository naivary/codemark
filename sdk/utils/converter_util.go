package utils

import (
	"reflect"
	"regexp"
	"strings"
)

const (
	TypeIDSep = "."
)

func Deref(typ reflect.Type) reflect.Type {
	if IsPointer(typ) {
		typ = typ.Elem()
	}
	return typ
}

func IsPointer(typ reflect.Type) bool {
	return typ.Kind() == reflect.Pointer
}

// ToType converts the value `v` to the Type `typ` handling pointer
// dereferencing and other inconveniences
func ToType(v reflect.Value, typ reflect.Type) (reflect.Value, error) {
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

func MatchTypeID(typeID, pattern string) bool {
	exp := regexp.MustCompile(pattern)
	return exp.MatchString(typeID)
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
	rtype := reflect.TypeOf(typ)
	return TypeID(rtype)
}

func typeID(typ reflect.Type, b *strings.Builder) string {
	if typ == nil {
		return b.String()
	}
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
