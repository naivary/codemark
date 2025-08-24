package generator

import (
	"fmt"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	outv1 "github.com/naivary/codemark/api/outputer/v1"
	"github.com/naivary/codemark/internal/config"
)

type name = string

type Manager struct {
	outputers map[name]outv1.Outputer
	cfg       map[string]any
}

func NewManager(configPath string) (*Manager, error) {
	const configSection = "outputer"
	mngr := &Manager{
		outputers: make(map[name]outv1.Outputer),
	}
	cfg, err := config.ReadIn(configPath, configSection)
	if err != nil {
		return nil, err
	}
	mngr.cfg = cfg
	return mngr, nil
}

func (m *Manager) Output(outputerName string, artifacts ...*genv1.Artifact) error {
	out, err := m.Get(outputerName)
	if err != nil {
		return err
	}
	cfg, isMap := m.cfg[outputerName].(map[string]any)
	if !isMap {
		cfg = make(map[string]any)
	}
	return out.Output(artifacts, cfg)
}

func (m *Manager) Get(name string) (outv1.Outputer, error) {
	out, found := m.outputers[name]
	if !found {
		return nil, fmt.Errorf("outputer not found: %s", name)
	}
	return out, nil
}

func (m *Manager) Add(out outv1.Outputer) error {
	name := out.Name()
	if _, found := m.outputers[name]; found {
		return fmt.Errorf("outputer already exists: %s", name)
	}
	m.outputers[name] = out
	return nil
}
