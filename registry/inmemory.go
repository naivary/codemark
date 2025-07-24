package registry

import (
	"fmt"
	"sync"

	"github.com/naivary/codemark/api/doc"
	optionapi "github.com/naivary/codemark/api/option"
)

func InMemory() Registry {
	return &inmem{
		opts: make(map[string]*optionapi.Option),
	}
}

var _ Registry = (*inmem)(nil)

type inmem struct {
	mu sync.Mutex

	opts map[string]*optionapi.Option
}

func (mem *inmem) Define(d *optionapi.Option) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	opt, exists := mem.opts[d.Ident]
	if exists {
		return fmt.Errorf("option already exists: %s", opt.Ident)
	}
	mem.opts[d.Ident] = d
	return nil
}

func (mem *inmem) Get(idn string) (*optionapi.Option, error) {
	opt, exists := mem.opts[idn]
	if exists {
		return opt, nil
	}
	return nil, fmt.Errorf("option not found: `%s`", idn)
}

func (mem *inmem) DocOf(ident string) (*doc.Option, error) {
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

func (mem *inmem) All() map[string]*optionapi.Option {
	return mem.opts
}
