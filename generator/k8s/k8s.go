package k8s

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/naivary/codemark/api"
	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/registry"
	"github.com/naivary/codemark/sdk"
	"golang.org/x/tools/go/packages"
)

var _ sdk.Generator = (*k8sGenerator)(nil)

type k8sGenerator struct {
	reg registry.Registry
}

func NewGenerator() (sdk.Generator, error) {
	reg, err := newRegistry()
	if err != nil {
		return nil, err
	}
	fmt.Println(reg.All())
	gen := &k8sGenerator{
		reg: reg,
	}
	return gen, nil
}

func (g k8sGenerator) Generate(infos map[*packages.Package]*loaderapi.Project) error {
	for _, proj := range infos {
		for _, strc := range proj.Structs {
			if shouldGenerateConfigMap(strc) {
				cm, err := createConfigMap(strc)
				if err != nil {
					return err
				}
				json.NewEncoder(os.Stdout).Encode(cm)
			}
		}
		for _, fn := range proj.Funcs {
			if shouldGeneratePod(fn) {
				pod, err := createPod(fn)
				if err != nil {
					return err
				}
				json.NewEncoder(os.Stdout).Encode(pod)
			}
		}
	}
	return nil
}

func (g k8sGenerator) Ressources() []string {
	return []string{
		"configmap",
		"pod",
		"meta",
	}
}

func (g k8sGenerator) Domain() string {
	return "k8s"
}

func (g k8sGenerator) Explain(ident string) api.OptionDoc {
	return api.OptionDoc{}
}

func (g k8sGenerator) OptionsOf(resource string) []api.OptionDoc {
	return nil
}

func (g k8sGenerator) Registry() registry.Registry {
	return g.reg
}
