package outputer

import (
	"fmt"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	outv1 "github.com/naivary/codemark/api/outputer/v1"
)

type name = string

type Manager struct {
	outputers map[name]outv1.Outputer
}

func NewManager(outputers ...outv1.Outputer) (*Manager, error) {
	mngr := &Manager{
		outputers: make(map[name]outv1.Outputer),
	}
	for _, outputer := range outputers {
		err := mngr.Add(outputer)
		if err != nil {
			return nil, err
		}
	}
	return mngr, nil
}

func (m *Manager) Output(name string, args []string, artifacts ...*genv1.Artifact) error {
	out, err := m.Get(name)
	if err != nil {
		return err
	}
	return out.Output(artifacts, args...)
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
