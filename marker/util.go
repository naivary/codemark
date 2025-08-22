package marker

import (
	"fmt"
	"go/types"
	"reflect"

	"github.com/naivary/codemark/rtypeutil"
	"github.com/naivary/codemark/typeutil"
)

// TypeOf returns the reflect.Type used for the given marker kind
func TypeOf(mkind Kind) reflect.Type {
	switch mkind {
	case STRING:
		return reflect.TypeFor[string]()
	case INT:
		return reflect.TypeFor[int64]()
	case FLOAT:
		return reflect.TypeFor[float64]()
	case COMPLEX:
		return reflect.TypeFor[complex128]()
	case BOOL:
		return reflect.TypeFor[bool]()
	case LIST:
		reflect.TypeFor[[]any]()
	}
	return nil
}

// KindFromRType returns which kind of marker the given type is. If no kind can be
// found INVALID will be returned.
func KindFromRType(typ reflect.Type) Kind {
	switch {
	case rtypeutil.IsValidSlice(typ):
		return LIST
	case rtypeutil.IsInt(typ), rtypeutil.IsUint(typ):
		return INT
	case rtypeutil.IsFloat(typ):
		return FLOAT
	case rtypeutil.IsComplex(typ):
		return COMPLEX
	case rtypeutil.IsBool(typ):
		return BOOL
	case rtypeutil.IsString(typ):
		return STRING
	}
	return INVALID
}

// IsTypedList is validating if all elements of `list` are of type `t`. If `t`
// is of kind `any` e.g. an empty interface nil will always be returned.
// Otherwise the following rules will be applied:
// 1. If `t` is array or slice the element type will be used for validation
// 2. If `t` is named type the underlying type will be used for validation iff
// one of the other rules apply
// 3. If `t` is a basic type the elements of the list will be validate for
// string, int*, float* and complex*
// 4. If `t` is an alias the elements will be validated for the reference type
// of the alias iff the type is fullfilling one of the other rules.
// 5. If `t` is any other type an error will be returned.
func IsTypedList(t types.Type, list []any) error {
	iface, isIface := t.(*types.Interface)
	if isIface && iface.Empty() {
		return nil
	}
	basic := typeutil.BasicTypeOf(t)
	if basic == nil {
		return fmt.Errorf("type is not fullfilling one of the rules described above: %v", t)
	}
	var err error
	switch basic.Kind() {
	case types.String:
		err = sameType[string](list)
	case types.Int, types.Int16, types.Int32, types.Int64:
		err = sameType[int64](list)
	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
		err = sameType[int64](list)
	case types.Float32, types.Float64:
		err = sameType[float64](list)
	case types.Complex64, types.Complex128:
		err = sameType[complex128](list)
	case types.Bool:
		err = sameType[bool](list)
	}
	return err
}

// sameType is validaignt if the elements of `list` are all of type `T`.
func sameType[T any](list []any) error {
	for _, el := range list {
		_, isT := el.(T)
		if !isT {
			return fmt.Errorf("%v is not of type %v", el, reflect.TypeFor[T]())
		}
	}
	return nil
}
