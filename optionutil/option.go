package optionutil

import (
	"fmt"
	"reflect"
	"strings"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
	"github.com/naivary/codemark/validate"
)

func Make(ident string, typ reflect.Type, doc *docv1.Option, isUnique bool, targets ...optionv1.Target) (optionv1.Option, error) {
	opt := optionv1.Option{
		Ident:   ident,
		Targets: targets,
		Type:    typ,
	}
	return opt, IsValid(opt)
}

func MustMake(ident string, output reflect.Type, doc *docv1.Option, isUnique bool, targets ...optionv1.Target) optionv1.Option {
	opt := optionv1.Option{
		Ident:   ident,
		Targets: targets,
		Type:    output,
	}
	if err := IsValid(opt); err != nil {
		panic(err)
	}
	return opt
}

func IsValid(opt optionv1.Option) error {
	if err := validate.Ident(opt.Ident); err != nil {
		return err
	}
	if opt.Type == nil {
		return fmt.Errorf("output type cannot be nil: %s", opt.Ident)
	}
	if len(opt.Targets) == 0 {
		return fmt.Errorf("definition has not target defined: %s", opt.Ident)
	}
	return nil
}

func DomainOf(ident string) string {
	s := strings.Split(ident, ":")
	if len(s) != 3 {
		return ""
	}
	return s[0]
}

func ResourceOf(ident string) string {
	s := strings.Split(ident, ":")
	if len(s) != 3 {
		return ""
	}
	return s[1]
}

func OptionOf(ident string) string {
	s := strings.Split(ident, ":")
	if len(s) != 3 {
		return ""
	}
	return s[2]
}
