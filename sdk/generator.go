package sdk

import (
	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/definition/target"
)

type OptionDoc struct {
	Targets []target.Target
	Doc     string
	Default string
}

type Generator interface {
	Domain() string
	Explain(ident string) OptionDoc
	// Ressources supported by this generator
	Ressources() []string
	Generate(proj *loaderapi.Project) ([]byte, error)
	OptionsOf(resource string) []OptionDoc
}
