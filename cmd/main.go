package main

import (
	"fmt"
	"os"

	"github.com/naivary/codemark"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

type required bool

func openAPIDefs() []*codemark.Definition {
	return []*codemark.Definition{
		codemark.MakeDef("openapi:validation:required", codemark.TargetField, required(false)),
	}
}

func run() error {
	reg := codemark.NewRegistry()
	for _, def := range openAPIDefs() {
		def.DeprecateInFavorOf("openapi:validation:required_more")
		if err := reg.Define(def); err != nil {
			return err
		}
	}
	conv, err := codemark.NewConverter(reg)
	if err != nil {
		return err
	}
	l := codemark.NewLoader(conv, nil)
	files, err := l.Load("./testdata")
	if err != nil {
		return err
	}
	_ = files
	return nil
}
