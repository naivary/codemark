package maker

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/parser/marker"
)

func MakeOption(idn string, output reflect.Type, targets ...core.Target) (*core.Option, error) {
	opt := &core.Option{
		Ident:   idn,
		Targets: targets,
		Output:  output,
	}
	return opt, opt.IsValid()
}

func MustMakeOpt(idn string, output reflect.Type, targets ...core.Target) *core.Option {
	opt := &core.Option{
		Ident:   idn,
		Targets: targets,
		Output:  output,
	}
	if err := opt.IsValid(); err != nil {
		panic(err)
	}
	return opt
}

func MakeOptWithDoc(name string, output reflect.Type, doc core.OptionDoc, targets ...core.Target) (*core.Option, error) {
	opt, err := MakeOption(name, output, targets...)
	if err != nil {
		return nil, err
	}
	opt.Doc = &doc
	return opt, opt.IsValid()
}

func MustMakeOptWithDoc(name string, output reflect.Type, doc core.OptionDoc, targets ...core.Target) *core.Option {
	opt, err := MakeOption(name, output, targets...)
	if err != nil {
		panic(err)
	}
	opt.Doc = &doc
	return opt
}

func MakeFakeMarker(mkind marker.Kind, value reflect.Value) marker.Marker {
	ident := fmt.Sprintf("codemark:fake:%s", mkind.String())
	return marker.New(ident, mkind, value)
}

func MakeMarker(ident string, mkind marker.Kind, value reflect.Value) (marker.Marker, error) {
	m := marker.New(ident, mkind, value)
	return m, m.IsValid()
}
