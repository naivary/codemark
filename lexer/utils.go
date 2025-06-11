package lexer

import (
	"fmt"
	"strings"
)

func IsValidIdent(ident string) error {
	plus := string(_plus)
	if strings.HasPrefix(ident, plus) {
		return fmt.Errorf("an identifier should not start with a plus. the plus is just like the `var` keyword in a regular programming language: %s\n", ident)
	}
	colon := string(_colon)
	numOfColons := strings.Count(ident, colon)
	if numOfColons < 2 {
		return fmt.Errorf("expected two colons in `%s` but got %d", ident, numOfColons)
	}
	for pathSegment := range strings.SplitSeq(ident, colon) {
		lastChar := rune(pathSegment[len(pathSegment)-1])
		if !isAlphaNumeric(lastChar) {
			return fmt.Errorf("identifier cannot end with an underscore `_` or dot `.`: %s\n", ident)
		}
	}
	return nil
}
