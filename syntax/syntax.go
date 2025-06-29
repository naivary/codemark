// syntax defines functions for shared rules between multiple types. For example
// the ident of a marker and definition must follow the same syntax otherwise
// the mapping is not possible.
package syntax

import (
	"fmt"
	"strings"
	"unicode"
)

// Ident validates if the given identifier is following the rule and conventions
// of codemark.
func Ident(ident string) error {
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
