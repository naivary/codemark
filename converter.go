package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
)

type Converter interface {
	Convert(marker parser.Marker) (reflect.Value, error)
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

func (c *converter) Convert(marker parser.Marker) (reflect.Value, error) {
	var empty reflect.Value
	// get definition from registry (done)
	// check which type the marker is and which type the definition output is
	//
	// if it's the same e.g. string and string then conert using
	// reflect.Value.Convert.
	//
	// if its not equal do a manual conversion by determining the kind
	name := marker.Ident()
	def := c.reg.Get(name)
	if def == nil {
		return empty, fmt.Errorf("marker `%s` is not defined in the registry", name)
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
	}

	return reflect.Value{}, nil
}

func convertString(marker parser.Marker, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	const runee = reflect.Int32
	const bytee = reflect.Uint8

	kind := def.Output.Kind()
	markerValue := marker.Value().String()
	value := reflect.New(def.Output).Elem()
	if kind == reflect.String {
		value.SetString(markerValue)
		return value, nil
	}
	if kind == bytee && len(markerValue) == 1 {
		b := byte(markerValue[0])
		v := reflect.ValueOf(b).Convert(def.Output)
		value.Set(v)
		return value, nil
	}
	if kind == runee && len(markerValue) == 1 {
		r := rune(markerValue[0])
		v := reflect.ValueOf(r).Convert(def.Output)
		value.Set(v)
		return value, nil
	}

	// check if its a []byte or []rune slice
	elemKind := def.Output.Elem().Kind()
	if kind == reflect.Slice && elemKind == bytee {
		bytes := []byte(markerValue)
		value.SetBytes(bytes)
		return value, nil
	}
	if kind == reflect.Slice && elemKind == runee {
		runes := []rune(markerValue)
		v := reflect.ValueOf(runes)
		value.Set(v)
		return value, nil
	}
	return empty, fmt.Errorf("cannot convert marker of kind `%v` to definition of kind `%v`", marker.Kind(), kind)
}

func convertBool(marker parser.Marker, def *Definition) (reflect.Value, error) {
	return marker.Value().Convert(def.Output), nil
}
