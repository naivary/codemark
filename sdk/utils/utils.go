package utils

import (
	"strings"
)

func Option(ident string) string {
	return strings.Split(ident, ":")[2]
}
