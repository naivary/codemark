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

func typeOfKind(k reflect.Kind) reflect.Type {
	switch k {
	case reflect.Int:
		return reflect.TypeOf(int(0))
	case reflect.Int8:
		return reflect.TypeOf(int8(0))
	case reflect.Int16:
		return reflect.TypeOf(int16(0))
	case reflect.Int32:
		return reflect.TypeOf(int32(0))
	case reflect.Int64:
		return reflect.TypeOf(int64(0))
	case reflect.Uint:
		return reflect.TypeOf(uint(0))
	case reflect.Uint8:
		return reflect.TypeOf(uint8(0))
	case reflect.Uint16:
		return reflect.TypeOf(uint16(0))
	case reflect.Uint32:
		return reflect.TypeOf(uint32(0))
	case reflect.Uint64:
		return reflect.TypeOf(uint64(0))
	case reflect.Float32:
		return reflect.TypeOf(float32(0.0))
	case reflect.Float64:
		return reflect.TypeOf(float64(0.0))
	case reflect.Complex64:
		return reflect.TypeOf(complex64(0 + 0i))
	case reflect.Complex128:
		return reflect.TypeOf(complex128(0 + 0i))
	case reflect.String:
		return reflect.TypeOf(string(""))
	case reflect.Bool:
		return reflect.TypeOf(bool(false))
	default:
		return nil
	}
}

func anyOf(k reflect.Kind, kinds ...reflect.Kind) bool {
	for _, kind := range kinds {
		if kind == k {
			return true
		}
	}
	return false
}

func typeID(typ reflect.Type, b *strings.Builder) (string, error) {
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
