package main

import "fmt"

type Registry interface {
	Define(def *Definition) error

	Get(name string) *Definition

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
	def, isDefined := r.defs[d.Name]
	if isDefined {
		return fmt.Errorf("definition is already defined: %s", def.Name)
	}
	r.defs[d.Name] = d
	return nil
}

func (r *registry) Get(name string) *Definition {
	return r.defs[name]
}

func (r *registry) All() map[string]*Definition {
	return r.defs
}
