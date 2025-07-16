package codemark

import (
	"github.com/naivary/codemark/generator"
)

// TODO: Mergging of registries is needed, otherwise its not possible to load
// the project
func GenerateWithManager(mngr *generator.Manager, pattern, domain string) error {
	gen, err := mngr.Get(domain)
	if err != nil {
		return err
	}
	reg := gen.Registry()
	infos, err := Load(reg, pattern)
	if err != nil {
		return err
	}
	return gen.Generate(infos)
}

func Generate(pattern, domain string) error {
	mngr, err := generator.NewManager()
	if err != nil {
		return err
	}
	return GenerateWithManager(mngr, pattern, domain)
}
