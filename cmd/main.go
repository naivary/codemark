package main

import (
	"fmt"
	"os"
)

func main() {
	if code, err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(code)
	}
}

func run() (int, error) {
	err := rootCmd.Execute()
	if err != nil {
		return 1, err
	}
	return 0, nil
}
