package utils

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/naivary/codemark/parser/marker"
)

func MarkerKindOf(typ reflect.Type) marker.Kind {
	kind := Deref(typ).Kind()
	switch kind {
	case reflect.Slice:
		return marker.LIST
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

func IsValidIdent(ident string) error {
	plus := "+"
	if strings.HasPrefix(ident, plus) {
		return fmt.Errorf("an identifier should not start with a plus. the plus is just like the `var` keyword in a regular programming language: %s\n", ident)
	}
	colon := ":"
	numOfColons := strings.Count(ident, colon)
	if numOfColons < 2 {
		return fmt.Errorf("expected two colons in `%s` but got %d", ident, numOfColons)
	}
	for pathSegment := range strings.SplitSeq(ident, colon) {
		lastChar := rune(pathSegment[len(pathSegment)-1])
		if !(unicode.IsLetter(lastChar) || unicode.IsDigit(lastChar)) {
			return fmt.Errorf("identifier cannot end with an underscore `_` or dot `.`: %s in %s\n", pathSegment, ident)
		}
	}
	return nil
}
