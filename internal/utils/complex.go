package utils

import (
	"fmt"
	"strings"
)

func ComplexOrder(c string) string {
	sign := string(c[0])
	cutSign := "+"
	if sign == "+" || sign == "-" {
		c = c[1:]
	} else {
		sign = "+"
	}
	x, y, found := strings.Cut(c, "+")
	if !found {
		x, y, _ = strings.Cut(c, "-")
		cutSign = "-"
	}
	if x[len(x)-1] == 'i' {
		return fmt.Sprintf("%s%s%s%s", cutSign, y, sign, x)
	}
	return fmt.Sprintf("%s%s%s%s", sign, x, cutSign, y)
}
