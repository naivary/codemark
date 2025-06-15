package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/naivary/codemark"
	"github.com/naivary/codemark/sdk"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

type required bool

type bools []bool

func openAPIDefs() []*sdk.Definition {
	return []*sdk.Definition{
		codemark.MustMakeDef("openapi_v3:validation:required", reflect.TypeFor[required](), sdk.TargetAny),
		codemark.MustMakeDef("openapi_v3:validation:bools", reflect.TypeFor[bools](), sdk.TargetAny),
	}
}

func run() error {
	reg := codemark.NewInMemoryRegistry()
	for _, def := range openAPIDefs() {
		if err := reg.Define(def); err != nil {
			return err
		}
	}
	conv, err := codemark.NewConvMngr(reg)
	if err != nil {
		return err
	}
	l := codemark.NewLocalLoader(conv, nil)
	proj, err := l.Load("./testdata")
	if err != nil {
		return err
	}
	fmt.Println(proj[0].Structs[0].Defs)
	return nil
}
