package codemark

import "reflect"

type ConverterType int

const (
	ConverterTypePtrToSlice = iota + 1
	ConverterTypePtrToArray
	ConverterTypeMap
	ConverterTypeString
	ConverterTypeInt
	ConverterTypeFloat
	ConverterTypeComplex
	ConverterTypeSlice
	ConverterTypeArr
	ConverterTypeNonEmptyStruct
)

func convTypeOf(def *Definition) ConverterType {
	isPtr := false
	kind := def.output.Kind()
	if kind == reflect.Pointer {
		isPtr = true
		kind = def.output.Elem().Kind()
	}
	if kind == reflect.Slice && isPtr {
		return ConverterTypePtrToSlice
	}
	if kind == reflect.Array && isPtr {
		return ConverterTypePtrToArray
	}
	if kind == reflect.Map {
		return ConverterTypeMap
	}
	if kind != reflect.Struct {
		return 0
	}
}
