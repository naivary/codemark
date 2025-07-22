package option

import (
	"reflect"

	"github.com/naivary/codemark/api/core"
)

func Make(idn string, output reflect.Type, targets ...core.Target) (*core.Option, error) {
	opt := &core.Option{
		Ident:   idn,
		Targets: targets,
		Output:  output,
	}
	return opt, opt.IsValid()
}

func MustMake(idn string, output reflect.Type, targets ...core.Target) *core.Option {
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

func MakeWithDoc(
	name string,
	output reflect.Type,
	doc core.OptionDoc,
	targets ...core.Target,
) (*core.Option, error) {
	opt, err := Make(name, output, targets...)
	if err != nil {
		return nil, err
	}
	opt.Doc = &doc
	return opt, opt.IsValid()
}

func MustMakeWithDoc(
	name string,
	output reflect.Type,
	doc core.OptionDoc,
	targets ...core.Target,
) *core.Option {
	opt, err := Make(name, output, targets...)
	if err != nil {
		panic(err)
	}
	opt.Doc = &doc
	return opt
}
