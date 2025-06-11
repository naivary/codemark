package codemark

import (
	"reflect"

	"github.com/naivary/codemark/sdk"
)

func MustMakeDef(idn string, t sdk.Target, output reflect.Type) *sdk.Definition {
	def := &sdk.Definition{
		Ident:  idn,
		Target: t,
		Output: output,
	}
	if err := def.IsValid(); err != nil {
		panic(err)
	}
	return def
}

func MustMakeDefWithHelp(name string, t sdk.Target, output reflect.Type, help *sdk.DefinitionHelp) *sdk.Definition {
	def := MustMakeDef(name, t, output)
	if err := def.IsValid(); err != nil {
		panic(err)
	}
	def.Help = help
	return def
}

func MakeDef(idn string, t sdk.Target, output reflect.Type) (*sdk.Definition, error) {
	def := &sdk.Definition{
		Ident:  idn,
		Target: t,
		Output: output,
	}
	return def, def.IsValid()
}

func MakeDefWithHelp(name string, t sdk.Target, output reflect.Type, help *sdk.DefinitionHelp) (*sdk.Definition, error) {
	def := MustMakeDef(name, t, output)
	if err := def.IsValid(); err != nil {
		panic(err)
	}
	def.Help = help
	return def, def.IsValid()
}
