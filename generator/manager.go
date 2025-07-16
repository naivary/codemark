package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

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

// Load is searching for the plugin and adds it to the manager
func (m *Manager) Load(name string) error {
	p := filepath.Join(m.dir, name)
	pl, err := plugin.Open(p)
	if err != nil {
		return err
	}
	sym, err := pl.Lookup("NewGenerator")
	if err != nil {
		return err
	}
	fn, ok := sym.(func() (sdk.Generator, error))
	if !ok {
		return fmt.Errorf("exported function NewGenerator of the plugin %s is not matching the interface func() (sdk.Generator, error)\n", name)
	}
	gen, err := fn()
	if err != nil {
		return err
	}
	return m.Add(gen)
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
