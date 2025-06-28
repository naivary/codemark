package maker

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk"
	"github.com/naivary/codemark/sdk/utils"
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

func MakeFakeDef(out reflect.Type) (*sdk.Definition, error) {
	ident := fmt.Sprintf("codemark:fake:%s", utils.NameFor(out))
	return MakeDef(ident, out, sdk.TargetAny)
}

func MakeFakeMarker(mkind marker.Kind, value reflect.Value) parser.Marker {
	ident := fmt.Sprintf("codemark:fake:%s", mkind.String())
	return parser.NewMarker(ident, mkind, value)
}
