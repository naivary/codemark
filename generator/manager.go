package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/naivary/codemark/sdk"
)

type domain = string

type Manager struct {
	dir  string
	gens map[domain]sdk.Generator
}

func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	mngr := &Manager{
		dir:  filepath.Join(home, ".codemark", "gens"),
		gens: make(map[domain]sdk.Generator),
	}
	return mngr, nil
}

func (m *Manager) Generate(domains ...string) error {
	return nil
}

func (m *Manager) Get(domain string) (sdk.Generator, error) {
	gen, found := m.gens[domain]
	if !found {
		return nil, fmt.Errorf("generator not found for domain: %s\n", domain)
	}
	return gen, nil
}

func (m *Manager) Add(gen sdk.Generator) error {
	domain := gen.Domain()
	if _, found := m.gens[domain]; found {
		return fmt.Errorf("generator for domain already exists: %s\n", domain)
	}
	m.gens[domain] = gen
	return nil
}
