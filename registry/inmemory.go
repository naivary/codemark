package registry

import (
	"fmt"
	"sync"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

func InMemory() Registry {
	return &inmem{
		opts: make(map[string]*optionv1.Option),
	}
}

var _ Registry = (*inmem)(nil)

type inmem struct {
	mu sync.Mutex

	opts map[string]*optionv1.Option
}

func (mem *inmem) Define(d *optionv1.Option) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	opt, exists := mem.opts[d.Ident]
	if exists {
		return fmt.Errorf("option already exists: %s", opt.Ident)
	}
	mem.opts[d.Ident] = d
	return nil
}

func (mem *inmem) Get(idn string) (*optionv1.Option, error) {
	opt, exists := mem.opts[idn]
	if exists {
		return opt, nil
	}
	return nil, fmt.Errorf("option not found: `%s`", idn)
}

func (mem *inmem) DocOf(ident string) (*docv1.Option, error) {
	opt, err := mem.Get(ident)
	if err != nil {
		return nil, err
	}
	return opt.Doc, nil
}

func (mem *inmem) Merge(reg Registry) error {
	for _, option := range reg.All() {
		if err := mem.Define(option); err != nil {
			return err
		}
	}
	return nil
}

func (mem *inmem) All() map[string]*optionv1.Option {
	return mem.opts
}
