package typeutil

import "go/types"

func BasicTypeOf(t types.Type) *types.Basic {
	basic, isBasic := t.(*types.Basic)
	if isBasic {
		return basic
	}
	switch t := t.(type) {
	case *types.Alias:
		return BasicTypeOf(t.Rhs())
	case *types.Slice:
		return BasicTypeOf(t.Elem())
	case *types.Array:
		return BasicTypeOf(t.Elem())
	case *types.Named:
		return BasicTypeOf(t.Underlying())
	case *types.Pointer:
		return BasicTypeOf(t.Elem())
	default:
		return nil
	}
}
