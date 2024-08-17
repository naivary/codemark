package main

import "reflect"

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
