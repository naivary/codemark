package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	if code, err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(code)
	}
}

func run(_ context.Context) (int, error) {
	return 0, nil
}
