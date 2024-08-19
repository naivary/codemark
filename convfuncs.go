package main

import (
	"errors"
	"fmt"
	"reflect"
)

func convComplex(c complex128, def *Definition) (any, error) {
	kind := def.kind
	if !isComplexInLimit(c, kind) {
		return nil, errors.New("does not fit")
	}
	if kind == reflect.Complex64 {
		return newComplexType[complex64](c, def), nil
	}
	if kind == reflect.Complex128 {
		return newComplexType[complex128](c, def), nil
	}
	return nil, errors.New("conversion not possible")
}

func convInt(i int64, def *Definition) (any, error) {
	kind := def.kind
	if !isIntInLimit(i, kind) {
		return nil, errors.New("is not fitting into the limit")
	}
	if kind == reflect.Int8 {
		return newNumberType[int8, int64](i, def), nil
	}
	if kind == reflect.Int16 {
		return newNumberType[int16, int64](i, def), nil
	}
	if kind == reflect.Int32 {
		return newNumberType[int32, int64](i, def), nil
	}
	if kind == reflect.Int64 {
		return newNumberType[int64, int64](i, def), nil
	}
	if kind == reflect.Int {
		return newNumberType[int, int64](i, def), nil
	}
	if kind == reflect.Float32 {
		return newNumberType[float32, int64](i, def), nil
	}
	if kind == reflect.Float64 {
		return newNumberType[float64, int64](i, def), nil
	}
	return nil, errors.New("conversion not possible")
}

func convUint(i uint64, def *Definition) (any, error) {
	kind := def.kind
	if !isUintInLimit(i, kind) {
		return nil, errors.New("integer is not fitting into it")
	}
	if kind == reflect.Uint8 {
		return newNumberType[uint8, uint64](i, def), nil
	}
	if kind == reflect.Uint16 {
		return newNumberType[uint16, uint64](i, def), nil
	}
	if kind == reflect.Uint32 {
		return newNumberType[uint32, uint64](i, def), nil
	}
	if kind == reflect.Uint64 {
		return newNumberType[uint64, uint64](i, def), nil
	}
	if kind == reflect.Uint {
		return newNumberType[uint, uint64](i, def), nil
	}
	return nil, errors.New("conversion no possible")
}

func convFloat(f float64, def *Definition) (any, error) {
	kind := def.kind
	if !isFloatInLimit(f, kind) {
		return nil, errors.New("does not fit in")
	}
	if kind == reflect.Float32 {
		return newNumberType[float32, float64](f, def), nil
	}
	if kind == reflect.Float64 {
		return newNumberType[float64, float64](f, def), nil
	}
	return nil, errors.New("conversion not possible")
}

func convString(s string, def *Definition) (any, error) {
	kind := def.kind
	if kind == reflect.String {
		return newType(s, def), nil
	}
	if kind == _byte && len(s) == 1 {
		return convByte(s[0], def)
	}
	if kind == _rune && len(s) == 1 {
		return convRune(rune(s[0]), def)
	}
	if kind == reflect.Slice && !anyOf(def.sliceKind, _rune, _byte) {
		return nil, fmt.Errorf("cannot convert it")
	}
	// check if its a []byte or []rune slice
	if def.sliceKind == _byte {
		bytes := []byte(s)
		return convSliceOfBytes(bytes, def)
	}
	if def.sliceKind == _rune {
		return convSliceOfRunes([]rune(s), def)
	}
	return nil, errors.New("conversion not possible")
}

func convByte(b byte, def *Definition) (any, error) {
	return newType(b, def), nil
}

func convRune(r rune, def *Definition) (any, error) {
	return newType(r, def), nil
}

func convSliceOfBytes(bytes []byte, def *Definition) (any, error) {
	return newType(bytes, def), nil
}

func convSliceOfRunes(runes []rune, def *Definition) (any, error) {
	return newType(runes, def), nil
}
