package converter

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/naivary/codemark"
	"github.com/naivary/codemark/marker"
)

type Converter interface {
	// IsPossible returns an error if the conversion of the `MarkerKind` is not
	// possible to `def.Output()`
	CanConvert(m marker.Marker, def *codemark.Definition) error

	//
	Convert(marker marker.Marker, def *codemark.Definition) (any, error)
}

func convString(s string, def *codemark.Definition) (reflect.Value, error) {
	var empty reflect.Value
	kind := def.kind
	if kind == reflect.Slice {
		kind = def.sliceKind()
	}
	if kind == reflect.String {
		return valueFor(s, def)
	}
	if kind == _byte && len(s) == 1 {
		b := byte(s[0])
		return valueFor(b, def)
	}
	if kind == _rune && len(s) == 1 {
		r := rune(s[0])
		return valueFor(r, def)
	}
	if kind == reflect.Interface {
		return valueFor(s, def)
	}
	return empty, fmt.Errorf("not convertiable to kind `%v`", kind)
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
	return toOutput(v, def)
}

func convBool(b bool, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	kind := def.kind
	if kind == reflect.Slice {
		kind = def.sliceKind()
	}
	if kind == reflect.Bool || kind == reflect.Interface {
		return valueFor(b, def)
	}
	return empty, fmt.Errorf("not convertiable to kind `%v`", def.kind)
}

func convertInteger(m marker.Marker, def *Definition) (any, error) {
	if !isIntConvPossible(def.kind) {
		return nil, errors.New("int conversion not possible")
	}
	var v reflect.Value
	var err error
	if isInt(def.kind) {
		v, err = convInt(m.Value().Int(), def)
	}
	if isUint(def.kind) {
		v, err = convUint(m.Value().Uint(), def)
	}
	if err != nil {
		return nil, err
	}
	return toOutput(v, def)
}

func convInt(i int64, def *Definition) (reflect.Value, error) {
	typ := def.nonPtrType()
	if typ.Kind() == reflect.Interface {
		return valueFor(i, def)
	}

	if typ.OverflowInt(i) {
		return reflect.Value{}, errors.New("overflow will occur")
	}
	return valueFor(i, def)
}

func convUint(i uint64, def *Definition) (reflect.Value, error) {
	typ := def.nonPtrType()
	if typ.Kind() == reflect.Interface {
		return valueFor(i, def)
	}
	if typ.OverflowUint(i) {
		return reflect.Value{}, errors.New("overflow will occur")
	}
	return valueFor(i, def)
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
	return toOutput(v, def)
}

func convFloat(f float64, def *Definition) (reflect.Value, error) {
	var empty reflect.Value
	typ := def.nonPtrType()
	if typ.Kind() == reflect.Interface {
		return valueFor(f, def)
	}
	if typ.OverflowFloat(f) {
		return empty, errors.New("overflow will occur")
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
	return toOutput(v, def)
}

func convComplex(c complex128, def *Definition) (reflect.Value, error) {
	typ := def.nonPtrType()
	if typ.Kind() == reflect.Interface {
		return valueFor(c, def)
	}
	if typ.OverflowComplex(c) {
		return reflect.Value{}, errors.New("overflow will occur")
	}
	return valueFor(c, def)
}

func convertSlice(m marker.Marker, def *Definition) (any, error) {
	if !isSliceConvPossible(def.kind) {
		return nil, errors.New("slice marker can only be converted to slice types")
	}
	elems := m.Value()
	slice := makeSlice(def.underlying, 0, elems.Len())
	if def.underlying.Kind() == reflect.Pointer {
		slice = reflect.New(slice.Type())
	}
	for i := range m.Value().Len() {
		elem := elems.Index(i)
		kind := elem.Elem().Kind()
		var v reflect.Value
		var err error
		if kind == reflect.String {
			v, err = convString(elem.Elem().String(), def)
		}
		if isInt(kind) {
			v, err = convInt(elem.Elem().Int(), def)
		}
		if isUint(kind) {
			v, err = convUint(elem.Elem().Uint(), def)
		}
		if kind == reflect.Bool {
			v, err = convBool(elem.Elem().Bool(), def)
		}
		if err != nil {
			return nil, err
		}
		slice = appendToSlice(slice, v)
	}
	return toOutput(slice, def)
}
