package registry

import (
	"fmt"
	"sync"

	"github.com/naivary/codemark/api"
)

func InMemory() Registry {
	return &inmem{
		defs: make(map[string]*api.Definition),
	}
}

var _ Registry = (*inmem)(nil)

type inmem struct {
	mu sync.Mutex

	defs map[string]*api.Definition
}

func (mem *inmem) Define(d *api.Definition) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	def, isDefined := mem.defs[d.Ident]
	if isDefined {
		return fmt.Errorf("definition is already defined: %s", def.Ident)
	}
	mem.defs[d.Ident] = d
	return nil
}

func (mem *inmem) Get(idn string) (*api.Definition, error) {
	def, found := mem.defs[idn]
	if found {
		return def, nil
	}
	return nil, fmt.Errorf("definition not found: `%s`", idn)
}

func (mem *inmem) DocOf(ident string) (*api.OptionDoc, error) {
	def, err := mem.Get(ident)
	if err != nil {
		return nil, err
	}
	return def.Doc, nil
}

func (mem *inmem) All() map[string]*api.Definition {
	return mem.defs
}
