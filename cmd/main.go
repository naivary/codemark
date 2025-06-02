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
		codemark.MakeDef("openapi_v3:validation:required", codemark.TargetField, required(false)),
	}
}

func run() error {
	reg := codemark.NewRegistry()
	for _, def := range openAPIDefs() {
		if err := reg.Define(def); err != nil {
			return err
		}
	}
	conv, err := codemark.NewConvMngr(reg)
	if err != nil {
		return err
	}
	l := codemark.NewLoader(conv, nil)
	files, err := l.Load("./testdata")
	if err != nil {
		return err
	}

	structs := files["github.com/naivary/codemark/testdata"][0].Structs
	for _, strc := range structs {
		for _, field := range strc.Fields {
			_ = field
		}
	}
	return nil
}
