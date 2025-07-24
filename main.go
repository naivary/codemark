package main

import (
	"os"

	"github.com/naivary/codemark/cmd"
)

func main() {
	code, err := cmd.Exec(nil, nil)
	if err != nil {
		os.Exit(code)
	}
}
