package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
)

type Converter interface {
	Convert(marker parser.Marker) (any, error)
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

func (c *converter) Convert(marker parser.Marker) (any, error) {
	name := marker.Ident()
	def := c.reg.Get(name)
	if def == nil {
		return nil, fmt.Errorf("marker `%s` is not defined in the registry", name)
	}
	switch marker.Kind() {
	// everything an be converted to any and the pointer of the type
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// can only be converted to float, string and int types and if it is
		// int32 to rune
	case reflect.Float64, reflect.Float32:
		// can only be conveted to  float types
	case reflect.Complex64, reflect.Complex128:
		// can only be converted to complext types
	case reflect.Bool:
		// only convertable to bool types
		return convertBool(marker, def)
	case reflect.String:
		return convertString(marker, def)
		// only convertable to string, []rune and []byte and if one letter rune
		// and byte too
	default:
		return nil, fmt.Errorf("invalid kind: `%s`", marker.Kind())
	}

	return reflect.Value{}, nil
}

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

func convertString(marker parser.Marker, def *Definition) (any, error) {
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

func convertBool(marker parser.Marker, def *Definition) (any, error) {
	value := reflect.New(def.Output).Elem()
	_, isPointer := ptrGuard(def)
	v := valueOf(marker.Value().Bool(), isPointer).Convert(def.Output)
	value.Set(v)
	return value.Interface(), nil
}
