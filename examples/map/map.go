package main

import (
	"log"
	"reflect"

	"github.com/naivary/codemark"
	"github.com/naivary/codemark/definition"
	"github.com/naivary/codemark/definition/target"
)

type Applier interface {
	ApplyTo(ident string)
}

type Map map[string]string

func (m Map) ApplyTo(ident string) {}

func defs() []*definition.Definition {
	return []*definition.Definition{
		codemark.MustMakeDef("openapi_v3:items:item.format", reflect.TypeFor[string](), target.FIELD),
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	reg := codemark.NewInMemRegistry()
	for _, def := range defs() {
		if err := reg.Define(def); err != nil {
			return err
		}
	}
	infos, err := codemark.Load(reg, "./project/")
	if err != nil {
		return err
	}
	_ = infos
	return nil
}
