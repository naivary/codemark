package option

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/naivary/codemark/api/doc"
	optionapi "github.com/naivary/codemark/api/option"
	"github.com/naivary/codemark/validate"
)

func IsValid(opt *optionapi.Option) error {
	if err := validate.Ident(opt.Ident); err != nil {
		return err
	}
	if opt.Output == nil {
		return fmt.Errorf("output type cannot be nil: %s", opt.Ident)
	}
	if len(opt.Targets) == 0 {
		return fmt.Errorf("definition has not target defined: %s", opt.Ident)
	}
	return nil
}

func Make(idn string, output reflect.Type, targets ...optionapi.Target) (*optionapi.Option, error) {
	opt := &optionapi.Option{
		Ident:   idn,
		Targets: targets,
		Output:  output,
	}
	return opt, IsValid(opt)
}

func MustMake(idn string, output reflect.Type, targets ...optionapi.Target) *optionapi.Option {
	opt := &optionapi.Option{
		Ident:   idn,
		Targets: targets,
		Output:  output,
	}
	if err := IsValid(opt); err != nil {
		panic(err)
	}
	return opt
}

func MakeWithDoc(
	name string,
	output reflect.Type,
	doc doc.Option,
	targets ...optionapi.Target,
) (*optionapi.Option, error) {
	opt, err := Make(name, output, targets...)
	if err != nil {
		return nil, err
	}
	opt.Doc = &doc
	return opt, IsValid(opt)
}

func MustMakeWithDoc(
	name string,
	output reflect.Type,
	doc doc.Option,
	targets ...optionapi.Target,
) *optionapi.Option {
	opt, err := Make(name, output, targets...)
	if err != nil {
		panic(err)
	}
	opt.Doc = &doc
	return opt
}

func DomainOf(ident string) string {
	s := strings.Split(ident, ":")
	if len(s) != 3 {
		return ""
	}
	return s[0]
}
