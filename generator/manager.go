package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

type domain = string

type Manager struct {
	dir  string
	gens map[domain]Generator
}

func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	mngr := &Manager{
		dir:  filepath.Join(home, ".codemark", "gens"),
		gens: make(map[domain]Generator),
	}
	return mngr, nil
}

func (m *Manager) Generate(domains ...string) error {
	return nil
}

func (m *Manager) Get(domain string) (Generator, error) {
	gen, found := m.gens[domain]
	if !found {
		return nil, fmt.Errorf("generator not found for domain: %s", domain)
	}
	return gen, nil
}

func (m *Manager) Add(gen Generator) error {
	domain := gen.Domain()
	if _, found := m.gens[domain]; found {
		return fmt.Errorf("generator for domain already exists: %s", domain)
	}
	m.gens[domain] = gen
	return nil
}
