package maker

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/definition"
	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk/utils"
)

func MakeDef(idn string, output reflect.Type, targets ...target.Target) (*definition.Definition, error) {
	def := &definition.Definition{
		Ident:   idn,
		Targets: targets,
		Output:  output,
	}
	return def, def.IsValid()
}

func MustMakeDef(idn string, output reflect.Type, targets ...target.Target) *definition.Definition {
	def := &definition.Definition{
		Ident:   idn,
		Targets: targets,
		Output:  output,
	}
	if err := def.IsValid(); err != nil {
		panic(err)
	}
	return def
}

func MakeDefWithHelp(name string, output reflect.Type, doc string, targets ...target.Target) (*definition.Definition, error) {
	def, err := MakeDef(name, output, targets...)
	if err != nil {
		return nil, err
	}
	def.Doc = doc
	return def, def.IsValid()
}

func MustMakeDefWithHelp(name string, output reflect.Type, doc string, targets ...target.Target) *definition.Definition {
	def, err := MakeDef(name, output, targets...)
	if err != nil {
		panic(err)
	}
	def.Doc = doc
	return def
}

func MakeFakeDef(out reflect.Type) (*definition.Definition, error) {
	ident := fmt.Sprintf("codemark:fake:%s", utils.NameFor(out))
	return MakeDef(ident, out, target.ANY)
}

func MakeFakeMarker(mkind marker.Kind, value reflect.Value) marker.Marker {
	ident := fmt.Sprintf("codemark:fake:%s", mkind.String())
	return marker.New(ident, mkind, value)
}

func MakeMarker(ident string, mkind marker.Kind, value reflect.Value) (marker.Marker, error) {
	m := marker.New(ident, mkind, value)
	return m, m.IsValid()
}
