package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/naivary/codemark/marker"
)

type Converter interface {
	// Convert is converting the given marker to the associated
	// `Definition.Output` type iff the target is correct and the conversion is
	// possible.
	Convert(marker marker.Marker, target Target) (any, error)
}

func NewConverter(reg Registry) (Converter, error) {
	if len(reg.All()) == 0 {
		return nil, errors.New("registry is empty")
	}
	m := &converter{
		reg: reg,
	}
	return m, nil
}

var _ Converter = (*converter)(nil)

type converter struct {
	reg Registry
}

func (c *converter) Convert(marker marker.Marker, target Target) (any, error) {
	name := marker.Ident()
	def := c.reg.Get(name)
	if def == nil {
		return nil, fmt.Errorf("marker `%s` is not defined in the registry", name)
	}
	if target != def.TargetType {
		return nil, fmt.Errorf("marker `%s` is appliable to `%s`. Was applied to `%s`", name, def.TargetType, target)
	}
	switch marker.Kind() {
	// everything an be converted to any and the pointer of the type
	case reflect.Int64:
		// can only be converted to float, string and int types and if it is
		// int32 to rune
		return convertNumber(marker, def)
	case reflect.Float64:
		// can only be conveted to  float types
		return convertFloat(marker, def)
	case reflect.Complex128:
		// can only be converted to complext types
	case reflect.Bool:
		// only convertable to bool types
		return convertBool(marker, def)
	case reflect.String:
		return convertString(marker, def)
		// only convertable to string, []rune and []byte and if one letter rune
		// and byte too
	}
	return nil, fmt.Errorf("invalid kind: `%s`", marker.Kind())
}

func convertUint(marker marker.Marker, def *Definition) (any, error) {
	i := uint64(marker.Value().Int())
	value := reflect.New(def.Output).Elem()
	kind, isPointer := ptrGuard(def)
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

func convertInt(marker marker.Marker, def *Definition) (any, error) {
	i := marker.Value().Int()
	value := reflect.New(def.Output).Elem()
	kind, isPointer := ptrGuard(def)
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
	if kind == reflect.Int && isIntInLimit(i, reflect.Int) {
		v := valueOf(int(i), isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	return nil, fmt.Errorf("cannot convert `%d` to `%v`. Be sure that the integer is fitting into the choosen bit size", i, def.Output)
}

func convertFloat(marker marker.Marker, def *Definition) (any, error) {
	return nil, nil
}

func convertNumber(marker marker.Marker, def *Definition) (any, error) {
	kind, _ := ptrGuard(def)
	i := marker.Value().Int()
	if i < 0 && isUint(kind) {
		return nil, fmt.Errorf("impossible uint conversion of `%d` to `%v`", i, def.Output)
	}
	if isUint(kind) {
		return convertUint(marker, def)
	}
	return convertInt(marker, def)
}

func convertString(marker marker.Marker, def *Definition) (any, error) {
	const runee = reflect.Int32
	const bytee = reflect.Uint8

	value := reflect.New(def.Output).Elem()
	markerValue := marker.Value().String()
	kind, isPointer := ptrGuard(def)
	if kind == reflect.String {
		v := valueOf(markerValue, isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	if kind == bytee && len(markerValue) == 1 {
		b := byte(markerValue[0])
		v := valueOf(b, isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}
	if kind == runee && len(markerValue) == 1 {
		r := rune(markerValue[0])
		v := valueOf(r, isPointer).Convert(def.Output)
		value.Set(v)
		return value.Interface(), nil
	}

	// check if its a []byte or []rune slice
	elemKind := def.Output.Elem().Kind()
	if kind == reflect.Slice && elemKind == bytee {
		bytes := []byte(markerValue)
		value.Set(valueOf(bytes, isPointer))
		return value.Interface(), nil
	}
	if kind == reflect.Slice && elemKind == runee {
		runes := []rune(markerValue)
		value.Set(valueOf(runes, isPointer))
		return value.Interface(), nil
	}
	return nil, fmt.Errorf("cannot convert marker of kind `%v` to definition of kind `%v`", marker.Kind(), kind)
}

func convertBool(marker marker.Marker, def *Definition) (any, error) {
	value := reflect.New(def.Output).Elem()
	_, isPointer := ptrGuard(def)
	v := valueOf(marker.Value().Bool(), isPointer).Convert(def.Output)
	value.Set(v)
	return value.Interface(), nil
}
