package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/naivary/codemark/marker"
)

type Converter interface {
	// Convert is converting the given marker to the associated
	// `Definition.output` type iff the target is correct and the conversion is
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
	if target != def.target {
		return nil, fmt.Errorf("marker `%s` is appliable to `%s`. Was applied to `%s`", name, def.target, target)
	}
	switch marker.Kind() {
	case reflect.Int64:
		return convertInteger(marker, def)
	case reflect.Float64:
		return convertDecimal(marker, def)
	case reflect.Complex128:
		return convertComplex(marker, def)
	case reflect.Bool:
		return convertBool(marker, def)
	case reflect.String:
		return convertString(marker, def)
	case reflect.Slice:
		return convertSlice(marker, def)
	}
	return nil, fmt.Errorf("invalid kind: `%s`", marker.Kind())
}

func convertString(m marker.Marker, def *Definition) (any, error) {
	if !isStringConvPossible(def.kind) {
		return nil, errors.New("string conversion not possible")
	}
	s := m.Value().String()
	value, err := convString(s, def)
	if err != nil {
		return nil, err
	}
	return convertToOutput(value, def)
}

func convString(s string, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	if def.kind == reflect.String {
		return valueFor(s, def)
	}
	if def.kind == _byte && len(s) == 1 {
		b := byte(s[0])
		return valueFor(b, def)
	}
	if def.kind == _rune && len(s) == 1 {
		r := rune(s[0])
		return valueFor(r, def)
	}
	return empty, fmt.Errorf("not convertiable to kind `%v`", def.kind)
}

func convertBool(m marker.Marker, def *Definition) (any, error) {
	if !isBoolConvPossible(def.kind) {
		return nil, errors.New("bool conversion not possible")
	}
	b := m.Value().Bool()
	v, err := convBool(b, def)
	if err != nil {
		return nil, err
	}
	return convertToOutput(v, def)
}

func convBool(b bool, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	if def.kind == reflect.Bool {
		return valueFor(b, def)
	}
	return empty, fmt.Errorf("not convertiable to kind `%v`", def.kind)
}

func convertInteger(m marker.Marker, def *Definition) (any, error) {
	if !isIntConvPossible(def.kind) {
		return nil, errors.New("int conversion not possible")
	}
	v, err := valueFor(m.Value().Interface(), def)
	if err != nil {
		return nil, err
	}
	return convertToOutput(v, def)
}

func convertDecimal(m marker.Marker, def *Definition) (any, error) {
	if !isFloatConvPossible(def.kind) {
		return nil, errors.New("float conversion not possible")
	}
	f := m.Value().Float()
	v, err := convFloat(f, def)
	if err != nil {
		return nil, err
	}
	return convertToOutput(v, def)
}

func convFloat(f float64, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	if !isFloatInLimit(f, def.kind) {
		return empty, errors.New("cannot convert overflow will occur")
	}
	return valueFor(f, def)
}

func convertComplex(m marker.Marker, def *Definition) (any, error) {
	if !isComplexConvPossible(def.kind) {
		return nil, errors.New("complex conversion not possible")
	}
	c := m.Value().Complex()
	v, err := convComplex(c, def)
	if err != nil {
		return nil, err
	}
	return convertToOutput(v, def)
}

func convComplex(c complex128, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	if !isComplexInLimit(c, def.kind) {
		return empty, errors.New("cannot convert overflow will occur")
	}
	return valueFor(c, def)
}

func convertSlice(m marker.Marker, def *Definition) (any, error) {
	if !isSliceConvPossible(def.kind) {
		return nil, errors.New("slice marker can only be converted to slice types")
	}
	elems := m.Value()
	slice := reflect.MakeSlice(def.underlying, 0, elems.Len())
	for i := 0; i < m.Value().Len(); i++ {
		elem := elems.Index(i)
		v, err := valueFor(elem.Interface(), def)
		if err != nil {
			return nil, err
		}
		slice = reflect.Append(slice, v)
	}
	return convertToOutput(slice, def)
}
