package codemark

import (
	"fmt"
	"reflect"
)

const (
	_rune = reflect.Int32
	_byte = reflect.Uint8
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

func isIntConvPossible(kind reflect.Kind) bool {
	return anyOf(kind,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
	)
}

func isSliceConvPossible(kind reflect.Kind) bool {
	return anyOf(kind, reflect.Slice)
}

func isFloatConvPossible(kind reflect.Kind) bool {
	return anyOf(kind, reflect.Float32, reflect.Float64)
}

func isComplexConvPossible(kind reflect.Kind) bool {
	return anyOf(kind, reflect.Complex64, reflect.Complex128)
}

func isBoolConvPossible(kind reflect.Kind) bool {
	return anyOf(kind, reflect.Bool, reflect.Struct)
}

func isStringConvPossible(kind reflect.Kind) bool {
	return anyOf(kind, reflect.String, _byte, _rune)
}

func isUint(k reflect.Kind) bool {
	return anyOf(k,
		reflect.Uint,
		reflect.Uint8, // e.g. byte
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
	)
}

func isInt(k reflect.Kind) bool {
	return anyOf(k,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32, // e.g. rune
		reflect.Int64,
	)
}

func valueFor[T any](val T, def *Definition) (reflect.Value, error) {
	typ := def.typ()
	// always call `Elem()` because typ is already the correct type.
	// Otherwise you might end with a double pointer
	decl := reflect.New(typ).Elem()
	if typ.Kind() == reflect.Pointer {
		return valueForPtr(val, def)
	}
	value := valueOf(val, def)
	if !value.CanConvert(typ) {
		return reflect.Value{}, fmt.Errorf("cannot convert from `%v` to `%v`", value.Type(), typ)
	}
	value = value.Convert(typ)
	decl.Set(value)
	return decl, nil
}

func valueForPtr[T any](val T, def *Definition) (reflect.Value, error) {
	typ := def.typ().Elem()
	value := valueOf(val, def).Elem()
	if !value.Type().ConvertibleTo(typ) {
		return reflect.Value{}, fmt.Errorf("type `%v` is not convertiable to `%v`", value.Type(), typ)
	}
	value = value.Convert(typ)
	decl := reflect.New(typ)
	decl.Elem().Set(value)
	return decl, nil
}

func isSliceElemPtr(def *Definition) bool {
	return def.kind == reflect.Slice && def.sliceType().Kind() == reflect.Pointer
}

func isNonSlicePtr(def *Definition) bool {
	return def.output.Kind() == reflect.Pointer && def.kind != reflect.Slice
}

func valueOf[T any](val T, def *Definition) reflect.Value {
	if isSliceElemPtr(def) {
		return reflect.ValueOf(&val)
	}
	if isNonSlicePtr(def) {
		return reflect.ValueOf(&val)
	}
	return reflect.ValueOf(val)
}

func toOutput(value reflect.Value, def *Definition) (any, error) {
	if !value.CanConvert(def.output) {
		return nil, fmt.Errorf("conversion from `%v` to `%v` is not possible", value.Type(), def.output)
	}
	output := value.Convert(def.output)
	return output.Interface(), nil
}

func underlying(t reflect.Type) reflect.Type {
	kind := t.Kind()
	if kind == reflect.Pointer {
		return reflect.PointerTo(t.Elem())
	}
	if kind == reflect.Slice {
		return reflect.SliceOf(t.Elem())
	}
	return typeOfKind(kind)
}

func makeSlice(t reflect.Type, l, c int) reflect.Value {
	kind := t.Kind()
	if kind != reflect.Pointer {
		return reflect.MakeSlice(t, l, c)
	}
	typ := t.Elem()
	return reflect.MakeSlice(typ, l, c)
}

func appendToSlice(slice reflect.Value, elem reflect.Value) reflect.Value {
	if slice.Kind() == reflect.Pointer {
		s := slice.Elem()
		s = reflect.Append(s, elem)
		slice.Elem().Set(s)
		return slice
	}
	return reflect.Append(slice, elem)
}
