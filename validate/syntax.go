package validate

import (
	"fmt"
	"strings"
	"unicode"
)

// Ident validates if the given identifier is following the rule and conventions
// of codemark. An identifier is composed of three components which are
// seperated by a comma (:) e.g. domain:resource:option. This convention MUST be
// followed because it will be used to load generators and make decisions.
func Ident(ident string) error {
	plus := "+"
	if strings.HasPrefix(ident, plus) {
		return fmt.Errorf(
			"an identifier should not start with a plus. the plus is just like the `var` keyword in a regular programming language: %s",
			ident,
		)
	}
	colon := ":"
	numOfColons := strings.Count(ident, colon)
	if numOfColons < 2 {
		return fmt.Errorf("expected two colons in `%s` but got %d", ident, numOfColons)
	}
	for pathSegment := range strings.SplitSeq(ident, colon) {
		lastChar := rune(pathSegment[len(pathSegment)-1])
		if !unicode.IsLetter(lastChar) && !unicode.IsDigit(lastChar) {
			return fmt.Errorf(
				"identifier cannot end with an underscore `_` or dot `.`: %s in %s",
				pathSegment,
				ident,
			)
		}
	}
	return nil
}
