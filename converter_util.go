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

func TypeID(typ reflect.Type, builder *strings.Builder) (string, error) {
	if typ == nil {
		return builder.String(), nil
	}
	if len(builder.String()) > 0 {
		builder.WriteString(".")
	}
	kind := typ.Kind()
	if _, err := builder.WriteString(kind.String()); err != nil {
		return _undefined, err
	}
	if kind == reflect.Pointer {
		return TypeID(typ.Elem(), builder)
	}
	if kind == reflect.Map {
		var keyBuilder strings.Builder
		key, _ := TypeID(typ.Key(), &keyBuilder)
		builder.WriteString(".")
		builder.WriteString(key)
		return TypeID(typ.Elem(), builder)
	}
	if kind == reflect.Slice {
		return TypeID(typ.Elem(), builder)
	}
	if kind == reflect.Chan {
		return TypeID(typ.Elem(), builder)
	}
	if kind == reflect.Array {
		return TypeID(typ.Elem(), builder)
	}
	if kind == reflect.Interface {
		return TypeID(nil, builder)
	}
	return TypeID(nil, builder)
}
