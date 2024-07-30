package main

type Registry interface {
	// Register is adding the definition to the
	// registry for potential later lookups
	Register(def *Definition) error

	// Lookup is returning the definition in the registry.
	// It's nil if it is not found
	Lookup(name string) *Definition

	// LookupHelp is returning the `DefinitionHelp` in the registry.
	// It's nil if it is not found
	LookupHelp(name string) *DefinitionHelp
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

func (r registry) Register(def *Definition) error {
	r.defs[def.Name] = def
	return nil
}

func (r registry) Lookup(name string) *Definition {
	return r.defs[name]
}

func (r registry) LookupHelp(name string) *DefinitionHelp {
	v, ok := r.defs[name]
	if !ok {
		return nil
	}
	return v.Help
}
