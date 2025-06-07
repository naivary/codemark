package codemark

import "fmt"

func NewInMemoryRegistry() Registry {
	return &inmemoryRegistry{
		defs: make(map[string]*Definition),
	}
}

var _ Registry = (*inmemoryRegistry)(nil)

type inmemoryRegistry struct {
	defs map[string]*Definition
}

func (mem *inmemoryRegistry) Define(d *Definition) error {
	def, isDefined := mem.defs[d.Ident]
	if isDefined {
		return fmt.Errorf("definition is already defined: %s", def.Ident)
	}
	mem.defs[d.Ident] = d
	return nil
}

func (mem *inmemoryRegistry) Get(idn string) (*Definition, error) {
	def, found := mem.defs[idn]
	if found {
		return def, nil
	}
	return nil, fmt.Errorf("definition not found: `%s`", idn)
}

func (mem *inmemoryRegistry) All() map[string]*Definition {
	return mem.defs
}
