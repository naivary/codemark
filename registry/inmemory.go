package registry

import (
	"fmt"
	"sync"

	"github.com/naivary/codemark/definition"
)

func InMemory() Registry {
	return &inmem{
		defs: make(map[string]*definition.Definition),
	}
}

var _ Registry = (*inmem)(nil)

type inmem struct {
	mu sync.Mutex

	defs map[string]*definition.Definition
}

func (mem *inmem) Define(d *definition.Definition) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	def, isDefined := mem.defs[d.Ident]
	if isDefined {
		return fmt.Errorf("definition is already defined: %s", def.Ident)
	}
	mem.defs[d.Ident] = d
	return nil
}

func (mem *inmem) Get(idn string) (*definition.Definition, error) {
	def, found := mem.defs[idn]
	if found {
		return def, nil
	}
	return nil, fmt.Errorf("definition not found: `%s`", idn)
}

func (mem *inmem) DocOf(ident string) (string, error) {
	def, err := mem.Get(ident)
	if err != nil {
		return "", nil
	}
	return def.Doc, nil
}

func (mem *inmem) All() map[string]*definition.Definition {
	return mem.defs
}
