package v1

import (
	"fmt"

	"github.com/naivary/codemark/optionutil"
)

type Options map[string][]any

func (o Options) Add(ident string, value any, isUnique bool) error {
	opts, ok := o[ident]
	if !ok {
		o[ident] = []any{value}
		return nil
	}
	if isUnique {
		return fmt.Errorf("unique option used more than once: %s", ident)
	}
	o[ident] = append(opts, value)
	return nil
}

func (o Options) Filter(domain, resource string) Options {
	opts := make(Options)
	for ident, opt := range o {
		if optionutil.ResourceOf(ident) == resource && optionutil.DomainOf(ident) == domain {
			opts[ident] = opt
		}
	}
	return opts
}

func (o Options) IsDefined(ident string) bool {
	return len(o[ident]) != 0
}

func (o Options) Get(ident string) ([]any, bool) {
	return o[ident], o.IsDefined(ident)
}
