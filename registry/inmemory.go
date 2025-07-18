package registry

import (
	"fmt"
	"sync"

	"github.com/naivary/codemark/api/core"
)

func InMemory() Registry {
	return &inmem{
		opts: make(map[string]*core.Option),
	}
}

var _ Registry = (*inmem)(nil)

type inmem struct {
	mu sync.Mutex

	opts map[string]*core.Option
}

func (mem *inmem) Define(d *core.Option) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	opt, exists := mem.opts[d.Ident]
	if exists {
		return fmt.Errorf("option already exists: %s", opt.Ident)
	}
	mem.opts[d.Ident] = d
	return nil
}

func (mem *inmem) Get(idn string) (*core.Option, error) {
	opt, exists := mem.opts[idn]
	if exists {
		return opt, nil
	}
	return nil, fmt.Errorf("option not found: `%s`", idn)
}

func (mem *inmem) DocOf(ident string) (*core.OptionDoc, error) {
	opt, err := mem.Get(ident)
	if err != nil {
		return nil, err
	}
	return opt.Doc, nil
}

func (mem *inmem) All() map[string]*core.Option {
	return mem.opts
}
