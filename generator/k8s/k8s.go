package k8s

import (
	"encoding/json"
	"os"

	"golang.org/x/tools/go/packages"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/registry"
)

var _ genv1.Generator = (*k8sGenerator)(nil)

type k8sGenerator struct {
	reg registry.Registry
}

func New() (genv1.Generator, error) {
	reg, err := newRegistry()
	if err != nil {
		return nil, err
	}
	gen := &k8sGenerator{
		reg: reg,
	}
	return gen, nil
}

func (g k8sGenerator) Domain() string {
	return "k8s"
}

func (g k8sGenerator) Explain(ident string) string {
	option, err := g.reg.Get(ident)
	if err != nil {
		return ""
	}
	return option.String()
}

func (g k8sGenerator) Registry() registry.Registry {
	return g.reg
}

func (g k8sGenerator) Generate(infos map[*packages.Package]*loaderapi.Information) error {
	file, err := os.Create("generate_manifests.yaml")
	if err != nil {
		return err
	}
	defer file.Close()
	for _, proj := range infos {
		for _, strc := range proj.Structs {
			if shouldGenerateConfigMap(strc) {
				cm, err := createConfigMap(strc)
				if err != nil {
					return err
				}
				err = json.NewEncoder(file).Encode(cm)
				if err != nil {
					return err
				}
			}
		}
		for _, fn := range proj.Funcs {
			if shouldGeneratePod(fn) {
				pod, err := createPod(fn)
				if err != nil {
					return err
				}
				err = json.NewEncoder(file).Encode(pod)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
