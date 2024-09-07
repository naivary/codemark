package main

import "fmt"

type Registry interface {
	Define(def *Definition) error

	Get(idn string) *Definition

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
	def, isDefined := r.defs[d.identifier]
	if isDefined {
		return fmt.Errorf("definition is already defined: %s", def.identifier)
	}
	r.defs[d.identifier] = d
	return nil
}

func (r *registry) Get(idn string) *Definition {
	return r.defs[idn]
}

func (r *registry) All() map[string]*Definition {
	return r.defs
}
