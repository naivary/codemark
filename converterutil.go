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

func isWithinLimit(i int64, limit reflect.Kind) bool {
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
