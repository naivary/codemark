package generator

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	regv1 "github.com/naivary/codemark/api/registry/v1"
	"github.com/naivary/codemark/loader"
	"github.com/naivary/codemark/registry"
)

func readInConfig() (map[string]any, error) {
	var config map[string]any
	file, err := os.ReadFile("codemark.yaml")
	if os.IsNotExist(err) {
		return map[string]any{}, nil
	}
	if err != nil {
		return nil, err
	}
	return config, yaml.Unmarshal(file, &config)
}

type domain = string

type Manager struct {
	gens map[domain]genv1.Generator
}

func NewManager() (*Manager, error) {
	mngr := &Manager{
		gens: make(map[domain]genv1.Generator),
	}
	return mngr, nil
}

func (m *Manager) Generate(pattern string, domains ...string) (map[domain][]*genv1.Artifact, error) {
	gens, err := m.all(domains...)
	if err != nil {
		return nil, err
	}
	reg, err := m.merge(gens...)
	if err != nil {
		return nil, err
	}
	info, err := loader.Load(reg, pattern)
	if err != nil {
		return nil, err
	}
	config, err := readInConfig()
	if err != nil {
		return nil, err
	}
	output := make(map[domain][]*genv1.Artifact)
	for _, gen := range gens {
		genCfg, isMap := config[gen.Domain()].(map[string]any)
		if !isMap {
			genCfg = make(map[string]any)
		}
		artifacts, err := gen.Generate(info, genCfg)
		if err != nil {
			return nil, err
		}
		output[gen.Domain()] = artifacts
	}
	return output, nil
}

func (m *Manager) Get(domain string) (genv1.Generator, error) {
	gen, found := m.gens[domain]
	if !found {
		return nil, fmt.Errorf("generator not found for domain: %s", domain)
	}
	return gen, nil
}

func (m *Manager) Add(gen genv1.Generator) error {
	domain := gen.Domain()
	if _, found := m.gens[domain]; found {
		return fmt.Errorf("generator for domain already exists: %s", domain)
	}
	m.gens[domain] = gen
	return nil
}

func (m *Manager) all(domains ...string) ([]genv1.Generator, error) {
	gens := make([]genv1.Generator, 0, len(domains))
	for _, domain := range domains {
		gen, err := m.Get(domain)
		if err != nil {
			return nil, err
		}
		gens = append(gens, gen)
	}
	return gens, nil
}

func (m *Manager) merge(gens ...genv1.Generator) (regv1.Registry, error) {
	regs := make([]regv1.Registry, 0, len(gens))
	for _, gen := range gens {
		regs = append(regs, gen.Registry())
	}
	return registry.Merge(regs...)
}
