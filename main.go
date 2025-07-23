package main

import (
	"fmt"
	"os"

	"github.com/naivary/codemark/cmd"
)

func main() {
	code, err := cmd.Exec(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(code)
	}
}
