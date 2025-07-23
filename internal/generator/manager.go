package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/loader"
	"github.com/naivary/codemark/registry"
)

type domain = string

type Manager struct {
	dir  string
	gens map[domain]generator.Generator
}

func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	mngr := &Manager{
		dir:  filepath.Join(home, ".codemark", "gens"),
		gens: make(map[domain]generator.Generator),
	}
	return mngr, nil
}

func (m *Manager) Generate(pattern string, domains ...string) error {
	gens := make([]generator.Generator, 0, len(domains))
	for _, domain := range domains {
		gen, err := m.Get(domain)
		if err != nil {
			return err
		}
		gens = append(gens, gen)
	}
	reg := registry.InMemory()
	for _, gen := range gens {
		if err := reg.Merge(gen.Registry()); err != nil {
			return err
		}
	}
	info, err := loader.Load(reg, pattern)
	if err != nil {
		return err
	}
	for _, gen := range gens {
		if err := gen.Generate(info); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) Get(domain string) (generator.Generator, error) {
	gen, found := m.gens[domain]
	if !found {
		return nil, fmt.Errorf("generator not found for domain: %s", domain)
	}
	return gen, nil
}

func (m *Manager) Add(gen generator.Generator) error {
	domain := gen.Domain()
	if _, found := m.gens[domain]; found {
		return fmt.Errorf("generator for domain already exists: %s", domain)
	}
	m.gens[domain] = gen
	return nil
}
