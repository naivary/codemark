package codemark

import (
	"reflect"
	"strings"
)

const (
	_rune = reflect.Int32
	_byte = reflect.Uint8

	_undefined = "UNDEFINED"
)

func typeID(typ reflect.Type, b *strings.Builder) (string, error) {
	if typ == nil {
		return b.String(), nil
	}
	if len(b.String()) > 0 {
		b.WriteString(".")
	}
	kind := typ.Kind()
	if _, err := b.WriteString(kind.String()); err != nil {
		return _undefined, err
	}
	if kind == reflect.Pointer {
		return typeID(typ.Elem(), b)
	}
	if kind == reflect.Map {
		var keyBuilder strings.Builder
		key, _ := typeID(typ.Key(), &keyBuilder)
		b.WriteString(".")
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
	if kind == reflect.Interface {
		return typeID(nil, b)
	}
	return typeID(nil, b)

}

func TypeID(typ reflect.Type) (string, error) {
	if typ == nil {
		return "", nil
	}
	var b strings.Builder
	return typeID(typ, &b)
}
