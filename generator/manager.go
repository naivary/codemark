package generator

import (
	"fmt"
	"maps"
	"slices"

	convv1 "github.com/naivary/codemark/api/converter/v1"
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

func NewManager(cfgFile string, gens ...genv1.Generator) (*Manager, error) {
	const configSection = "gens"
	mngr := &Manager{
		gens: make(map[domain]genv1.Generator),
	}
	cfg, err := config.ReadIn(cfgFile, configSection)
	if err != nil {
		return nil, err
	}
	mngr.cfg = cfg
	for _, gen := range gens {
		err := mngr.Add(gen)
		if err != nil {
			return nil, err
		}
	}
	return mngr, nil
}

func (m *Manager) Domains() []string {
	return slices.Collect(maps.Keys(m.gens))
}

func (m *Manager) All() map[domain]genv1.Generator {
	return m.gens
}

func (m *Manager) Generate(convs []convv1.Converter, pattern string) (map[domain][]*genv1.Artifact, error) {
	reg, err := m.merge(m.allGens())
	if err != nil {
		return nil, err
	}
	info, err := loader.Load(reg, convs, pattern)
	if err != nil {
		return nil, err
	}
	output := make(map[domain][]*genv1.Artifact)
	for _, gen := range m.gens {
		artifacts, err := gen.Generate(info, m.configFor(gen))
		if err != nil {
			return nil, err
		}
		output[gen.Domain().Name] = artifacts
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
	domain := gen.Domain().Name
	if _, found := m.gens[domain]; found {
		return fmt.Errorf("generator for domain already exists: %s", domain)
	}
	m.gens[domain] = gen
	return nil
}

func (m *Manager) allGens() []genv1.Generator {
	return slices.Collect(maps.Values(m.gens))
}

func (m *Manager) configFor(gen genv1.Generator) map[string]any {
	genCfg, isMap := m.cfg[gen.Domain().Name].(map[string]any)
	if isMap {
		return genCfg
	}
	return make(map[string]any)
}

func (m *Manager) merge(gens []genv1.Generator) (regv1.Registry, error) {
	regs := make([]regv1.Registry, 0, len(gens))
	for _, gen := range gens {
		regs = append(regs, gen.Registry())
	}
	return registry.Merge(regs...)
}
