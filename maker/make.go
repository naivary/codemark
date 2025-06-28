package maker

import (
	"reflect"

	"github.com/naivary/codemark/sdk"
)

func MakeDef(idn string, output reflect.Type, targets ...sdk.Target) (*sdk.Definition, error) {
	def := &sdk.Definition{
		Ident:   idn,
		Targets: targets,
		Output:  output,
	}
	return def, def.IsValid()
}

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

func MakeDefWithHelp(name string, output reflect.Type, help *sdk.DefinitionHelp, targets ...sdk.Target) (*sdk.Definition, error) {
	def, err := MakeDef(name, output, targets...)
	if err != nil {
		return nil, err
	}
	def.Help = help
	return def, def.IsValid()
}

func MustMakeDefWithHelp(name string, output reflect.Type, help *sdk.DefinitionHelp, targets ...sdk.Target) *sdk.Definition {
	def, err := MakeDef(name, output, targets...)
	if err != nil {
		panic(err)
	}
	def.Help = help
	return def
}
