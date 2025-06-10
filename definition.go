package codemark

import (
	"reflect"

	"github.com/naivary/codemark/sdk"
)

func MakeDef(idn string, t sdk.Target, output reflect.Type) *sdk.Definition {
	def := &sdk.Definition{
		Ident:  idn,
		Target: t,
		Output: output,
	}
	return def
}

func MakeDefWithHelp(name string, t sdk.Target, output reflect.Type, help *sdk.DefinitionHelp) *sdk.Definition {
	def := MakeDef(name, t, output)
	def.Help = help
	return def
}
