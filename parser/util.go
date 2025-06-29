package parser

import (
	"strings"
)

// complexOrder returns the imaginary part as the first return value and the
// real part as the second from a string.
func complexOrder(s string) (string, string) {
	x, y, _ := strings.Cut(s, "+")
	if x[len(x)-1] == 'i' {
		return x, y
	}
	return y, x
}
