package main

import (
	"math"
	"reflect"
)

func valueOf[T any](v T, isPointer bool) reflect.Value {
	if isPointer {
		return reflect.ValueOf(&v)
	}
	return reflect.ValueOf(v)
}

func ptrGuard(def *Definition) (reflect.Kind, bool) {
	kind := def.Output.Kind()
	isPointer := false
	if kind == reflect.Ptr {
		isPointer = true
		kind = def.Output.Elem().Kind()
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
		return f <= math.MaxFloat32 && f >= -math.SmallestNonzeroFloat32
	case reflect.Float64:
		return f <= math.MaxFloat64 && f >= -math.SmallestNonzeroFloat64
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
