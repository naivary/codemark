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
	}
	return nil, fmt.Errorf("invalid kind: `%s`", marker.Kind())
}

func convertString(m marker.Marker, def *Definition) (any, error) {
	if !isStringConvPossible(def) {
		return nil, errors.New("string conversion not possible")
	}
	s := m.Value().String()
	value, err := convString(s, def)
	if err != nil {
		return nil, err
	}
	return value.Interface(), nil
}

func convString(s string, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	if def.kind == reflect.String {
		return newValue(s, def)
	}
	if def.kind == _byte && len(s) == 1 {
		b := byte(s[0])
		return newValue(b, def)
	}

	if def.kind == _rune && len(s) == 1 {
		r := rune(s[0])
		return newValue(r, def)
	}
	return empty, fmt.Errorf("not convertiable to kind `%v`", def.kind)
}

func convertBool(m marker.Marker, def *Definition) (any, error) {
	if !isBoolConvPossible(def) {
		return nil, errors.New("bool conversion not possible")
	}
	b := m.Value().Bool()
	v, err := convBool(b, def)
    if err != nil {
        return nil, err
    }
	return v.Interface(), nil
}

func convBool(b bool, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	if def.kind == reflect.Bool {
		return newValue(b, def)
	}
	return empty, fmt.Errorf("not convertiable to kind `%v`", def.kind)
}

func convertInteger(m marker.Marker, def *Definition) (any, error) {
	var v reflect.Value
	var err error
	if !isIntConvPossible(def) {
		return nil, errors.New("int conversion not possible")
	}
	if isUint(def.kind) {
		v, err = convUint(m.Value().Uint(), def)
	}
	if isInt(def.kind) {
		v, err = convInt(m.Value().Int(), def)
	}
    if err != nil {
        return nil, err
    }
	return v.Interface(), nil
}

func convInt(i int64, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	if !isIntInLimit(i, def.kind) {
		return empty, errors.New("cannot convert overflow will occur")
	}
	return newValue(i, def)
}

func convUint(i uint64, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	if !isUintInLimit(i, def.kind) {
		return empty, errors.New("cannot convert overflow will occur")
	}
	return newValue(i, def)
}

func convertDecimal(m marker.Marker, def *Definition) (any, error) {
	if !isFloatConvPossible(def) {
		return nil, errors.New("float conversion not possible")
	}
	f := m.Value().Float()
	v, err := convFloat(f, def)
    if err != nil {
        return nil, err
    }
	return v.Interface(), nil
}

func convFloat(f float64, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	if !isFloatInLimit(f, def.kind) {
		return empty, errors.New("cannot convert overflow will occur")
	}
	return newValue(f, def)
}

func convertComplex(m marker.Marker, def *Definition) (any, error) {
	if !isComplexConvPossible(def) {
		return nil, errors.New("complex conversion not possible")
	}
	c := m.Value().Complex()
	v, err := convComplex(c, def)
    if err != nil {
        return nil, err
    }
	return v.Interface(), nil
}

func convComplex(c complex128, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	if !isComplexInLimit(c, def.kind) {
		return empty, errors.New("cannot convert overflow will occur")
	}
	return newValue(c, def)
}
