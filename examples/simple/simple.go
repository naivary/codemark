package main

import (
	"log"
	"reflect"

	"github.com/naivary/codemark"
	"github.com/naivary/codemark/definition"
	"github.com/naivary/codemark/definition/target"
)

type Applier interface {
	ApplyTo(s *Schema)
}

type Schema struct {
	Required bool
}

type required bool

func (r required) ApplyTo(s *Schema) {
	s.Required = bool(r)
}

func defs() []*definition.Definition {
	return []*definition.Definition{
		codemark.MustMakeDef("openapi_v3:validation:required", reflect.TypeFor[required](), target.FIELD),
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
	sch := &Schema{}
	for _, proj := range infos {
		for _, s := range proj.Structs {
			for _, field := range s.Fields {
				for _, defs := range field.Defs {
					for _, def := range defs {
						marker := def.(Applier)
						marker.ApplyTo(sch)
					}
				}
			}
		}
	}
	return nil
}
