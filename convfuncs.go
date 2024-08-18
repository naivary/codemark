package main

import (
	"fmt"
	"reflect"
)

func convComplex(c complex128, def *Definition) (any, error) {
	kind, isPointer := resolvePtr(def)
	value := reflect.New(def.Output).Elem()
	if kind != reflect.Complex64 && kind != reflect.Complex128 {
		return nil, fmt.Errorf("conversion of `%f+%fi` is not possible to a non complex type `%v`", real(c), imag(c), def.Output)
	}
	if kind == reflect.Complex64 && isComplexInLimit(c, reflect.Complex64) {
		v := valueOf(complex64(c), isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	v := valueOf(c, isPointer).Convert(def.Output)
	value.Set(v)
	return value.Interface(), nil
}

func convInt(i int64, def *Definition) (any, error) {
	value := reflect.New(def.Output).Elem()
	kind, isPointer := resolvePtr(def)
	if kind == reflect.Int8 && isIntInLimit(i, reflect.Int8) {
		v := valueOf(int8(i), isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	if kind == reflect.Int16 && isIntInLimit(i, reflect.Int16) {
		v := valueOf(int16(i), isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	if kind == reflect.Int32 && isIntInLimit(i, reflect.Int32) {
		v := valueOf(int32(i), isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	if kind == reflect.Int64 && isIntInLimit(i, reflect.Int64) {
		v := valueOf(i, isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	v := valueOf(int(i), isPointer).Convert(def.Output)
	value.Set(v)
	return value.Interface(), nil
}

func convUint(i uint64, def *Definition) (any, error) {
	value := reflect.New(def.Output).Elem()
	kind, isPointer := resolvePtr(def)
	if kind == reflect.Uint8 && isUintInLimit(i, reflect.Uint8) {
		v := valueOf(uint8(i), isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	if kind == reflect.Uint16 && isUintInLimit(i, reflect.Uint16) {
		v := valueOf(uint16(i), isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	if kind == reflect.Uint32 && isUintInLimit(i, reflect.Uint32) {
		v := valueOf(uint32(i), isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	if kind == reflect.Uint64 && isUintInLimit(i, reflect.Uint64) {
		v := valueOf(i, isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	if kind == reflect.Uint && isUintInLimit(i, reflect.Uint) {
		v := valueOf(uint(i), isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	return nil, fmt.Errorf("cannot convert `%d` to `%v`", i, def.Output)
}

func convFloat(f float64, def *Definition) (any, error) {
	kind, isPointer := resolvePtr(def)
	value := reflect.New(def.Output).Elem()
	if kind == reflect.Float32 && isFloatInLimit(f, reflect.Float32) {
		v := valueOf(float32(f), isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	v := valueOf(f, isPointer).Convert(def.Output)
	value.Set(v)
	return value.Interface(), nil
}

func convString(s string, def *Definition) (any, error) {
	const runee = reflect.Int32
	const bytee = reflect.Uint8

	value := reflect.New(def.Output).Elem()
	kind, isPointer := resolvePtr(def)
	if kind == reflect.String {
		v := valueOf(s, isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	if kind == bytee && len(s) == 1 {
		b := byte(s[0])
		v := valueOf(b, isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	if kind == runee && len(s) == 1 {
		r := rune(s[0])
		v := valueOf(r, isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}

	// check if its a []byte or []rune slice
	elemKind := def.Output.Elem().Kind()
	if kind == reflect.Slice && elemKind == bytee {
		bytes := []byte(s)
		value.Set(valueOf(bytes, isPointer))
		return value.Interface(), nil
	}
	runes := []rune(s)
	value.Set(valueOf(runes, isPointer))
	return value.Interface(), nil
}
