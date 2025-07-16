package generator

import (
	"fmt"

	"github.com/naivary/codemark/sdk"
)

type domain = string

type Manager struct {
	gens map[domain]sdk.Generator
}

func NewManager() (*Manager, error) {
	mngr := &Manager{
		gens: make(map[domain]sdk.Generator),
	}
	return mngr, nil
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
