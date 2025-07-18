package maker

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/parser/marker"
)

func MakeOption(idn string, output reflect.Type, targets ...core.Target) (*core.Option, error) {
	def := &core.Option{
		Ident:   idn,
		Targets: targets,
		Output:  output,
	}
	return def, def.IsValid()
}

func MustMakeOpt(idn string, output reflect.Type, targets ...core.Target) *core.Option {
	def := &core.Option{
		Ident:   idn,
		Targets: targets,
		Output:  output,
	}
	if err := def.IsValid(); err != nil {
		panic(err)
	}
	return def
}

func MakeOptWithDoc(name string, output reflect.Type, doc core.OptionDoc, targets ...core.Target) (*core.Option, error) {
	def, err := MakeOption(name, output, targets...)
	if err != nil {
		return nil, err
	}
	def.Doc = &doc
	return def, def.IsValid()
}

func MustMakeOptWithDoc(name string, output reflect.Type, doc core.OptionDoc, targets ...core.Target) *core.Option {
	def, err := MakeOption(name, output, targets...)
	if err != nil {
		panic(err)
	}
	def.Doc = &doc
	return def
}

func MakeFakeMarker(mkind marker.Kind, value reflect.Value) marker.Marker {
	ident := fmt.Sprintf("codemark:fake:%s", mkind.String())
	return marker.New(ident, mkind, value)
}

func MakeMarker(ident string, mkind marker.Kind, value reflect.Value) (marker.Marker, error) {
	m := marker.New(ident, mkind, value)
	return m, m.IsValid()
}
