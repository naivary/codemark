package k8s

import (
	"encoding/json"
	"os"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	loaderv1 "github.com/naivary/codemark/api/loader/v1"
	"github.com/naivary/codemark/registry"
)

var _ genv1.Generator = (*generator)(nil)

type generator struct {
	reg registry.Registry
}

func New() (genv1.Generator, error) {
	reg, err := newRegistry()
	if err != nil {
		return nil, err
	}
	gen := &generator{
		reg: reg,
	}
	return gen, nil
}

func (g generator) Domain() string {
	return "k8s"
}

func (g generator) Explain(ident string) string {
	option, err := g.reg.Get(ident)
	if err != nil {
		return ""
	}
	return option.String()
}

func (g generator) Registry() registry.Registry {
	return g.reg
}

func (g generator) Generate(proj loaderv1.Project) ([]genv1.Artifact, error) {
	for _, info := range proj {
		for _, strc := range info.Structs {
			if shouldGenerateConfigMap(strc) {
				cm, err := createConfigMap(strc)
				if err != nil {
					return nil, err
				}
				err = json.NewEncoder(os.Stdout).Encode(cm)
				if err != nil {
					return nil, err
				}
			}
		}
		for _, fn := range info.Funcs {
			if isMainFunc(fn) {
				pod, err := createPod(fn)
				if err != nil {
					return nil, err
				}
				role, err := createRBACRole(fn)
				if err != nil {
					return nil, err
				}
				err = json.NewEncoder(os.Stdout).Encode(pod)
				err = json.NewEncoder(os.Stdout).Encode(role)
			}
		}
	}
	return nil, nil
}
