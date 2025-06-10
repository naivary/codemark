package codemark

import (
	"fmt"

	"github.com/naivary/codemark/sdk"
)

func NewInMemoryRegistry() sdk.Registry {
	return &inmemoryRegistry{
		defs: make(map[string]*sdk.Definition),
	}
}

var _ sdk.Registry = (*inmemoryRegistry)(nil)

type inmemoryRegistry struct {
	defs map[string]*sdk.Definition
}

func (mem *inmemoryRegistry) Define(d *sdk.Definition) error {
	def, isDefined := mem.defs[d.Ident]
	if isDefined {
		return fmt.Errorf("definition is already defined: %s", def.Ident)
	}
	mem.defs[d.Ident] = d
	return nil
}

func (mem *inmemoryRegistry) Get(idn string) (*sdk.Definition, error) {
	def, found := mem.defs[idn]
	if found {
		return def, nil
	}
	return nil, fmt.Errorf("definition not found: `%s`", idn)
}

func (mem *inmemoryRegistry) All() map[string]*sdk.Definition {
	return mem.defs
}
