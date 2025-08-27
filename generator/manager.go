package generator

import (
	"fmt"
	"maps"
	"slices"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	regv1 "github.com/naivary/codemark/api/registry/v1"
	"github.com/naivary/codemark/internal/config"
	"github.com/naivary/codemark/loader"
	"github.com/naivary/codemark/registry"
)

type domain = string

type Manager struct {
	gens map[domain]genv1.Generator

	cfg map[string]any
}

func NewManager(configPath string) (*Manager, error) {
	const configSection = "gens"
	mngr := &Manager{
		gens: make(map[domain]genv1.Generator),
	}
	cfg, err := config.ReadIn(configPath, configSection)
	if err != nil {
		return nil, err
	}
	mngr.cfg = cfg
	return mngr, nil
}

func (m *Manager) Generate(pattern string) (map[domain][]*genv1.Artifact, error) {
	reg, err := m.merge(slices.Collect(maps.Values(m.gens)))
	if err != nil {
		return nil, err
	}
	info, err := loader.Load(reg, pattern)
	if err != nil {
		return nil, err
	}
	output := make(map[domain][]*genv1.Artifact)
	for _, gen := range m.gens {
		genCfg, isMap := m.cfg[gen.Domain()].(map[string]any)
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

func (m *Manager) merge(gens []genv1.Generator) (regv1.Registry, error) {
	regs := make([]regv1.Registry, 0, len(gens))
	for _, gen := range gens {
		regs = append(regs, gen.Registry())
	}
	return registry.Merge(regs...)
}
