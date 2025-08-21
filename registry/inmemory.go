package registry

import (
	"errors"
	"fmt"
	"sync"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
	regv1 "github.com/naivary/codemark/api/registry/v1"
)

var ErrRegistryEmpty = errors.New(
	"the registry contains no definitions. Make sure to add your definitions using the `Define` method",
)

func InMemory() regv1.Registry {
	return &inmem{
		opts: make(map[string]*optionv1.Option),
	}
}

var _ regv1.Registry = (*inmem)(nil)

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

func (mem *inmem) All() map[string]*optionv1.Option {
	return mem.opts
}
