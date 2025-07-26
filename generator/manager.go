package generator

import (
	"fmt"
	"os"
	"path/filepath"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	"github.com/naivary/codemark/loader"
	"github.com/naivary/codemark/registry"
)

type domain = string

type Manager struct {
	dir  string
	gens map[domain]genv1.Generator
}

func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	mngr := &Manager{
		dir:  filepath.Join(home, ".codemark", "gens"),
		gens: make(map[domain]genv1.Generator),
	}
	return mngr, nil
}

func (m *Manager) Generate(pattern string, domains ...string) (map[domain][]*genv1.Artifact, error) {
	gens := make([]genv1.Generator, 0, len(domains))
	for _, domain := range domains {
		gen, err := m.Get(domain)
		if err != nil {
			return nil, err
		}
		gens = append(gens, gen)
	}
	reg := registry.InMemory()
	for _, gen := range gens {
		if err := reg.Merge(gen.Registry()); err != nil {
			return nil, err
		}
	}
	info, err := loader.Load(reg, pattern)
	if err != nil {
		return nil, err
	}
	artifacts := make(map[domain][]*genv1.Artifact, len(gens))
	for _, gen := range gens {
		generatedArtifacts, err := gen.Generate(info)
		if err != nil {
			return nil, err
		}
		artifacts[gen.Domain()] = generatedArtifacts
	}
	return artifacts, nil
}

func (m *Manager) Get(domain string) (genv1.Generator, error) {
	gen, found := m.gens[domain]
	if !found {
		return nil, fmt.Errorf("genv1.not found for domain: %s", domain)
	}
	return gen, nil
}

func (m *Manager) Add(gen genv1.Generator) error {
	domain := gen.Domain()
	if _, found := m.gens[domain]; found {
		return fmt.Errorf("genv1.for domain already exists: %s", domain)
	}
	m.gens[domain] = gen
	return nil
}
