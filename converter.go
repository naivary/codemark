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
		return convertNumber(marker, def)
	case reflect.Float64:
		return convertFloat(marker, def)
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

func convertComplex(marker marker.Marker, def *Definition) (any, error) {
	c := marker.Value().Complex()
	if !isComplexConvPossible(def) {
		return nil, errImpossibleConv(marker, def)
	}
	return convComplex(c, def)
}

func convertFloat(marker marker.Marker, def *Definition) (any, error) {
	f := marker.Value().Float()
	if !anyOf(def.kind, reflect.Float32, reflect.Float64) {
		return nil, fmt.Errorf("cannot convert `%f` to `%v`. Conversion from float to integer will not be handled", f, def.output)
	}
	return convFloat(f, def)
}

func convertNumber(marker marker.Marker, def *Definition) (any, error) {
	i := marker.Value().Int()
	if i < 0 && isUint(def.kind) {
		return nil, fmt.Errorf("impossible uint conversion of `%d` to `%v`", i, def.output)
	}
	if isUint(def.kind) {
		return convUint(uint64(i), def)
	}
	return convInt(i, def)
}

func convertString(m marker.Marker, def *Definition) (any, error) {
	if !isStringConvPossible(def) {
		return nil, errImpossibleConv(m, def)
	}
	s := m.Value().String()
	return convString(s, def)
}

func convertBool(marker marker.Marker, def *Definition) (any, error) {
	value := reflect.New(def.output).Elem()
	v := valueOf(marker.Value().Bool(), def.isPointer).Convert(def.output)
	value.Set(v)
	return value.Interface(), nil
}

func convertSlice(m marker.Marker, def *Definition) (any, error) {
	slice := m.Value()
	if def.kind != reflect.Slice {
		return nil, errors.New("cannot convert slice to a non-slice kind")
	}
    // TODO: has to do the same conversion rules for these as for the single
    // elements
	switch def.sliceKind {
	case reflect.Interface:
		return newSliceType[any](slice, def)
	case reflect.String:
		return newSliceType[string](slice, def)
	case reflect.Int:
		return newSliceType[int](slice, def)
	case reflect.Int8:
		return newSliceType[int8](slice, def)
	case reflect.Int16:
		return newSliceType[int16](slice, def)
	case reflect.Int32:
		return newSliceType[int32](slice, def)
	case reflect.Int64:
		return newSliceType[int64](slice, def)
	case reflect.Uint:
		return newSliceType[uint](slice, def)
	case reflect.Uint8:
		return newSliceType[uint8](slice, def)
	case reflect.Uint16:
		return newSliceType[uint16](slice, def)
	case reflect.Uint32:
		return newSliceType[uint32](slice, def)
	case reflect.Uint64:
		return newSliceType[uint64](slice, def)
	case reflect.Float32:
		return newSliceType[float32](slice, def)
	case reflect.Float64:
		return newSliceType[float64](slice, def)
	case reflect.Complex64:
		return newSliceType[complex64](slice, def)
	case reflect.Complex128:
		return newSliceType[complex128](slice, def)
	default:
		return nil, errors.New("cannot convert it")
	}
}
