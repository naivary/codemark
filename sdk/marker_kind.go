//go:generate stringer -type=MarkerKind

package sdk

import (
	"reflect"

	sdkutil "github.com/naivary/codemark/sdk/utils"
)

type MarkerKind int

const (
	MarkerKindString MarkerKind = iota + 1
	MarkerKindFloat
	MarkerKindInt
	MarkerKindComplex
	MarkerKindBool
	MarkerKindList
)

func MarkerKindOf(typ reflect.Type) MarkerKind {
	kind := sdkutil.Deref(typ).Kind()
	switch kind {
	case reflect.Slice:
		return MarkerKindList
	// rune=int32 & byte=uint8
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return MarkerKindInt
	case reflect.Float32, reflect.Float64:
		return MarkerKindFloat
	case reflect.Complex64, reflect.Complex128:
		return MarkerKindComplex
	case reflect.Bool:
		return MarkerKindBool
	case reflect.String:
		return MarkerKindString
	}
	return 0
}
