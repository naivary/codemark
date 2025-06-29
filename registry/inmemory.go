package registry

import (
	"fmt"
	"sync"

	"github.com/naivary/codemark/sdk"
)

func InMemory() sdk.Registry {
	return &inmemoryRegistry{
		defs: make(map[string]*sdk.Definition),
	}
}

var _ sdk.Registry = (*inmemoryRegistry)(nil)

type inmemoryRegistry struct {
	mu sync.Mutex

	defs map[string]*sdk.Definition
}

func (mem *inmemoryRegistry) Define(d *sdk.Definition) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	def, isDefined := mem.defs[d.Ident]
	if isDefined {
		return fmt.Errorf("definition is already defined: %s", def.Ident)
	}
	mem.defs[d.Ident] = d
	return nil
}

func (mem *inmemoryRegistry) Get(idn string) (*sdk.Definition, error) {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	def, found := mem.defs[idn]
	if found {
		return def, nil
	}
	return nil, fmt.Errorf("definition not found: `%s`", idn)
}

func (mem *inmemoryRegistry) All() map[string]*sdk.Definition {
	return mem.defs
}
