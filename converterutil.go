package main

import (
	"errors"
	"math"
	"reflect"
)

const (
	_rune = reflect.Int32
	_byte = reflect.Uint8
)

func isConvertibleTo(from reflect.Kind, to reflect.Kind) bool {
	if from == reflect.Slice {
		return isSliceConvPossible(to)
	}
	if from == reflect.String {
		return isStringConvPossible(to)
	}
	if from == reflect.Int64 {
		return isIntConvPossible(to)
	}
	if from == reflect.Float64 {
		return isFloatConvPossible(to)
	}
	if from == reflect.Bool {
		return isBoolConvPossible(to)
	}
	if from == reflect.Complex128 {
		return isComplexConvPossible(to)
	}
	return false
}

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
	return anyOf(kind, reflect.String, reflect.Uint8, reflect.Int32)
}

func isIntInLimit(i int64, limit reflect.Kind) bool {
	switch limit {
	case reflect.Int:
		return i <= math.MaxInt && i >= math.MinInt
	case reflect.Int8:
		return i <= math.MaxInt8 && i >= math.MinInt8
	case reflect.Int16:
		return i <= math.MaxInt16 && i >= math.MinInt16
	case reflect.Int32:
		return i <= math.MaxInt32 && i >= math.MinInt32
	case reflect.Int64:
		return i <= math.MaxInt64 && i >= math.MinInt64
	default:
		return false
	}
}

func isUintInLimit(i uint64, limit reflect.Kind) bool {
	switch limit {
	case reflect.Uint:
		return i <= math.MaxUint
	case reflect.Uint8:
		return i <= math.MaxUint8
	case reflect.Uint16:
		return i <= math.MaxUint16
	case reflect.Uint32:
		return i <= math.MaxUint32
	case reflect.Uint64:
		return i <= math.MaxUint64
	default:
		return false
	}
}

func isFloatInLimit(f float64, limit reflect.Kind) bool {
	switch limit {
	case reflect.Float32:
		return f <= math.MaxFloat32 && f >= -math.MaxFloat32
	case reflect.Float64:
		return f <= math.MaxFloat64 && f >= -math.MaxFloat64
	default:
		return false
	}
}

func isComplexInLimit(c complex128, limit reflect.Kind) bool {
	r := real(c)
	img := imag(c)
	switch limit {
	case reflect.Complex64:
		return isFloatInLimit(r, reflect.Float32) && isFloatInLimit(img, reflect.Float32)
	case reflect.Complex128:
		return isFloatInLimit(r, reflect.Float64) && isFloatInLimit(img, reflect.Float64)
	default:
		return false
	}
}

func isUint(k reflect.Kind) bool {
	return anyOf(k, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64)
}

func isInt(k reflect.Kind) bool {
	return anyOf(k, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64)
}

func valueFor[T any](val T, def *Definition) (reflect.Value, error) {
	typ := def.Type()
	// always call `Elem()` because typ is already the correct type.
	// Otherwise you might end with a double pointer
	decl := reflect.New(typ).Elem()
	if typ.Kind() == reflect.Ptr {
		return valueForPtr(val, def)
	}
	value := valueOf(val, def)
	if !value.CanConvert(typ) {
		return reflect.Value{}, errors.New("cannot convert")
	}
	value = value.Convert(typ)
	decl.Set(value)
	return decl, nil
}

func valueForPtr[T any](val T, def *Definition) (reflect.Value, error) {
	typ := def.Type().Elem()
	value := valueOf(val, def).Elem()
	if !value.Type().ConvertibleTo(typ) {
		return reflect.Value{}, errors.New("types are not convertiable")
	}
	value = value.Convert(typ)
	decl := reflect.New(typ)
	decl.Elem().Set(value)
	return decl, nil
}

func isElemOfSlicePtr(def *Definition) bool {
	return def.kind == reflect.Slice && def.sliceType().Kind() == reflect.Ptr
}

func isNonSlicePtr(def *Definition) bool {
	return def.output.Kind() == reflect.Ptr && def.kind != reflect.Slice
}

func valueOf[T any](val T, def *Definition) reflect.Value {
	if isElemOfSlicePtr(def) {
		return reflect.ValueOf(&val)
	}
	if isNonSlicePtr(def) {
		return reflect.ValueOf(&val)
	}
	return reflect.ValueOf(val)
}

func convertToOutput(value reflect.Value, def *Definition) (any, error) {
	if !value.CanConvert(def.output) {
		return nil, errors.New("cannot convert")
	}
	output := value.Convert(def.output)
	return output.Interface(), nil
}

func underlying(t reflect.Type) reflect.Type {
	kind := t.Kind()
	if kind == reflect.Ptr {
		return reflect.PointerTo(t.Elem())
	}
	if kind == reflect.Slice {
		return reflect.SliceOf(t.Elem())
	}
	return typeOfKind(kind)
}

func makeSlice(t reflect.Type, l, c int) reflect.Value {
	kind := t.Kind()
	if kind != reflect.Ptr {
		return reflect.MakeSlice(t, l, c)
	}
	typ := t.Elem()
	return reflect.MakeSlice(typ, l, c)
}
