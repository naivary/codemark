package codemark

import (
	"fmt"
)

type Registry interface {
	Define(def *Definition) error

	Get(idn string) (*Definition, error)

	All() map[string]*Definition
}

func NewRegistry() Registry {
	return &registry{
		defs: make(map[string]*Definition),
	}
}

var _ Registry = (*registry)(nil)

type registry struct {
	defs map[string]*Definition
}

func (r *registry) Define(d *Definition) error {
	def, isDefined := r.defs[d.Ident]
	if isDefined {
		return fmt.Errorf("definition is already defined: %s", def.Ident)
	}
	r.defs[d.Ident] = d
	return nil
}

func (r *registry) Get(idn string) (*Definition, error) {
	def, found := r.defs[idn]
	if found {
		return def, nil
	}
	return nil, fmt.Errorf("definition not found: `%s`", idn)
}

func (r *registry) All() map[string]*Definition {
	return r.defs
}
