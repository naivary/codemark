package k8s

import (
	"encoding/json"
	"os"

	"golang.org/x/tools/go/packages"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	loaderapi "github.com/naivary/codemark/api/loader"
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

func (g generator) Generate(infos map[*packages.Package]*loaderapi.Information) error {
	for _, proj := range infos {
		for _, strc := range proj.Structs {
			if shouldGenerateConfigMap(strc) {
				cm, err := createConfigMap(strc)
				if err != nil {
					return err
				}
				err = json.NewEncoder(os.Stdout).Encode(cm)
				if err != nil {
					return err
				}
			}
		}
		for _, fn := range proj.Funcs {
			if isMainFunc(fn) {
				pod, err := createPod(fn)
				if err != nil {
					return err
				}
				role, err := createRBACRole(fn)
				if err != nil {
					return err
				}
				err = json.NewEncoder(os.Stdout).Encode(pod)
				err = json.NewEncoder(os.Stdout).Encode(role)
			}
		}
	}
	return nil
}
