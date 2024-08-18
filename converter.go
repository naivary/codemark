package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/naivary/codemark/marker"
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
	kind, _ := resolvePtr(def)
	return anyOf(kind,
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
	kind, _ := resolvePtr(def)
	return anyOf(kind, reflect.Float32, reflect.Float64)
}

func isComplexConvPossible(def *Definition) bool {
	kind, _ := resolvePtr(def)
	return anyOf(kind, reflect.Complex64, reflect.Complex128)
}

func isBoolConvPossible(def *Definition) bool {
	kind, _ := resolvePtr(def)
	return kind == reflect.Bool
}

func isStringConvPossible(def *Definition) bool {
	kind, _ := resolvePtr(def)
	if anyOf(kind, reflect.String, reflect.Uint8, reflect.Int32) {
		return true
	}
	if kind != reflect.Slice {
		return false
	}
	sliceKind := def.Output.Elem().Kind()
	return anyOf(sliceKind, _byte, _rune)
}

func errImpossibleConv(m marker.Marker, def *Definition) error {
	kind, _ := resolvePtr(def)
	return fmt.Errorf("cannot conver marker of kind `%v` to definition of kind `%v`", m.Kind(), kind)
}

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
	//TODO: everything an be converted to any and the pointer of the type
	case reflect.Int64:
		return convertNumber(marker, def)
	case reflect.Float64:
		return convertFloat(marker, def)
	case reflect.Complex128:
		return convertComplex(marker, def)
	case reflect.Bool:
		return convertBool(marker, def)
	case reflect.String:
		return convertString(marker, def)
	}
	return nil, fmt.Errorf("invalid kind: `%s`", marker.Kind())
}

func convertComplex(marker marker.Marker, def *Definition) (any, error) {
	c := marker.Value().Complex()
	if !isComplexConvPossible(def) {
		return nil, errImpossibleConv(marker, def)
	}
	return convComplex(c, def)
}

func convertFloat(marker marker.Marker, def *Definition) (any, error) {
	kind, _ := resolvePtr(def)
	f := marker.Value().Float()
	if kind != reflect.Float32 && kind != reflect.Float64 {
		return nil, fmt.Errorf("cannot convert `%f` to `%v`. Conversion from float to integer will not be handled", f, def.Output)
	}
	return convFloat(f, def)
}

func convertNumber(marker marker.Marker, def *Definition) (any, error) {
	kind, _ := resolvePtr(def)
	i := marker.Value().Int()
	if i < 0 && isUint(kind) {
		return nil, fmt.Errorf("impossible uint conversion of `%d` to `%v`", i, def.Output)
	}
	if isUint(kind) && isUintInLimit(uint64(i), kind) {
		return convUint(uint64(i), def)
	}
	if isIntInLimit(i, kind) {
		return convInt(i, def)
	}
	return nil, errImpossibleConv(marker, def)
}

func convertString(m marker.Marker, def *Definition) (any, error) {
	if !isStringConvPossible(def) {
		return nil, errImpossibleConv(m, def)
	}
	s := m.Value().String()
	return convString(s, def)
}

func convertBool(marker marker.Marker, def *Definition) (any, error) {
	value := reflect.New(def.Output).Elem()
	_, isPointer := resolvePtr(def)
	v := valueOf(marker.Value().Bool(), isPointer).Convert(def.Output)
	value.Set(v)
	return value.Interface(), nil
}
