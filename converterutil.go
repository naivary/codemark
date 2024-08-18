package main

import (
	"math"
	"reflect"
)

const (
	_rune = reflect.Int32
	_byte = reflect.Uint8
)

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
		reflect.Int32, // +rune
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8, // +byte
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
	)
}

func isFloatConvPossible(def *Definition) bool {
	return anyOf(def.kind, reflect.Float32, reflect.Float64)
}

func isComplexConvPossible(def *Definition) bool {
	return anyOf(def.kind, reflect.Complex64, reflect.Complex128)
}

func isBoolConvPossible(def *Definition) bool {
	return def.kind == reflect.Bool
}

func isStringConvPossible(def *Definition) bool {
	if anyOf(def.kind, reflect.String, reflect.Uint8, reflect.Int32) {
		return true
	}
	if def.kind != reflect.Slice {
		return false
	}
	sliceKind := def.output.Elem().Kind()
	return anyOf(sliceKind, _byte, _rune)
}

func valueOf[T any](v T, isPointer bool) reflect.Value {
	if isPointer {
		return reflect.ValueOf(&v)
	}
	return reflect.ValueOf(v)
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
	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}
