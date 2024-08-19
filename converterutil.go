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

func isIntConvPossible(def *Definition) bool {
	return anyOf(def.kind,
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

func isFloatConvPossible(def *Definition) bool {
	return anyOf(def.kind, reflect.Float32, reflect.Float64)
}

func isComplexConvPossible(def *Definition) bool {
	return anyOf(def.kind, reflect.Complex64, reflect.Complex128)
}

func isBoolConvPossible(def *Definition) bool {
	return anyOf(def.kind, reflect.Bool, reflect.Struct)
}

func isStringConvPossible(def *Definition) bool {
	return anyOf(def.kind, reflect.String, reflect.Uint8, reflect.Int32)
}

func resolvePtr(def *Definition) (reflect.Kind, bool) {
	kind := def.output.Kind()
	isPointer := false
	if kind == reflect.Ptr {
		isPointer = true
		kind = def.output.Elem().Kind()
	}
	return kind, isPointer
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

func newValue[T any](val T, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	kindValue, err := convertToKind(val, def)
	if err != nil {
		return empty, err
	}
	return convertToDef(kindValue, def)
}

func convertToDef(kindValue reflect.Value, def *Definition) (reflect.Value, error) {
	defValue := reflect.New(def.output).Elem()
	converted := kindValue.Convert(def.output)
	defValue.Set(converted)
	return defValue, nil
}

func convertToKind[T any](val T, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	typ := typeOfKind(def.kind)
	convertTo := reflect.New(typ)
    if !def.isPointer {
        convertTo = convertTo.Elem()
    }
	value := valueOf(val, def.isPointer)
	if !value.CanConvert(convertTo.Type()) {
		return empty, errors.New("cannot convert")
	}
	return value.Convert(convertTo.Type()), nil
}

func valueOf[T any](v T, isPointer bool) reflect.Value {
	if isPointer {
		return reflect.ValueOf(&v)
	}
	return reflect.ValueOf(v)
}
