package codemark

import (
	"reflect"

	"github.com/naivary/codemark/sdk"
)

func MustMakeDef(idn string, output reflect.Type, targets ...sdk.Target) *sdk.Definition {
	def := &sdk.Definition{
		Ident:   idn,
		Targets: targets,
		Output:  output,
	}
	if err := def.IsValid(); err != nil {
		panic(err)
	}
	return def
}

func MustMakeDefWithHelp(name string, output reflect.Type, help *sdk.DefinitionHelp, targets ...sdk.Target) *sdk.Definition {
	def := MustMakeDef(name, output, targets...)
	if err := def.IsValid(); err != nil {
		panic(err)
	}
	def.Help = help
	return def
}

func MakeDef(idn string, output reflect.Type, targets ...sdk.Target) (*sdk.Definition, error) {
	def := &sdk.Definition{
		Ident:   idn,
		Targets: targets,
		Output:  output,
	}
	return def, def.IsValid()
}

func MakeDefWithHelp(name string, output reflect.Type, help *sdk.DefinitionHelp, targets ...sdk.Target) (*sdk.Definition, error) {
	def := MustMakeDef(name, output, targets...)
	if err := def.IsValid(); err != nil {
		panic(err)
	}
	def.Help = help
	return def, def.IsValid()
}
