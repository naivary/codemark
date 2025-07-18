package core

import (
	"bytes"
	"fmt"
)

type Docer[T any] interface {
	Doc() T
}

// TODO: make this a custom package
func trunc(s string, n int) string {
	var b bytes.Buffer
	pos := 1
	for _, r := range s {
		if pos%n == 0 && r == ' ' {
			fmt.Fprintf(&b, "\n")
			pos = 1
			continue
		}
		if r == '\n' {
			pos = 1
		}
		if pos%n != 0 {
			pos++
		}
		fmt.Fprint(&b, string(r))
	}
	return b.String()
}

type OptionDoc struct {
	Doc     string
	Default string
}
