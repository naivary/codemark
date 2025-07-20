package marker

import "reflect"

// TypeForMarkerKind returns the reflect.Type used for the given marker kind
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

func KindOf(typ reflect.Type) Kind {
	kind := Deref(typ).Kind()
	switch kind {
	case reflect.Slice:
		return LIST
	// rune=int32 & byte=uint8
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return marker.INT
	case reflect.Float32, reflect.Float64:
		return marker.FLOAT
	case reflect.Complex64, reflect.Complex128:
		return marker.COMPLEX
	case reflect.Bool:
		return marker.BOOL
	case reflect.String:
		return marker.STRING
	}
	return 0
}
