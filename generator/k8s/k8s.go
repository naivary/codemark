package k8s

import (
	"slices"

	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/registry"
	"github.com/naivary/codemark/sdk"
	"golang.org/x/tools/go/packages"
)

var _ sdk.Generator = (*k8sGenerator)(nil)

type k8sGenerator struct {
	reg registry.Registry
}

func New() (sdk.Generator, error) {
	reg := registry.InMemory()
	if err := provisionRegistry(reg); err != nil {
		return nil, err
	}
	gen := &k8sGenerator{
		reg: reg,
	}
	return gen, nil
}

func (g k8sGenerator) Generate(infos map[*packages.Package]*loaderapi.Project) error {
	for _, proj := range infos {
		for _, strc := range proj.Structs {
			for _, field := range strc.Fields {
				for ident := range field.Defs {
					if ident == "k8s:configmap:default" {
						createConfigMap(strc)
					}
				}
			}
		}
	}
	return nil
}

func (g k8sGenerator) Ressources() []string {
	return []string{
		"configmap",
	}
}

func (g k8sGenerator) Domain() string {
	return "k8s"
}

func (g k8sGenerator) Explain(ident string) sdk.OptionDoc {
	return sdk.OptionDoc{}
}

func (g k8sGenerator) OptionsOf(resource string) []sdk.OptionDoc {
	return nil
}

func (g k8sGenerator) Registry() registry.Registry {
	return g.reg
}

func provisionRegistry(reg registry.Registry) error {
	defs := slices.Concat(configMapDefs())
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
		}
	}
	return nil
}
