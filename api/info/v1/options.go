package v1

import (
	"fmt"
)

type Options map[string][]any

func (o Options) Add(idn string, value any, isUnique bool) error {
	opts, ok := o[idn]
	if !ok {
		o[idn] = []any{value}
		return nil
	}
	if isUnique {
		return fmt.Errorf("option is unique but was used more than once: %s", idn)
	}
	o[idn] = append(opts, value)
	return nil
}

func (o Options) Filter(resource string) {}

func (o Options) IsOptDefined(ident string) {}

func (o Options) Len(ident string) int {
	return len(o[ident])
}
