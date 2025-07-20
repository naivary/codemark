package generator

import (
	"github.com/naivary/codemark/api/core"
	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/registry"
	"golang.org/x/tools/go/packages"
)

type Generator interface {
	// Domain for which the generator is responsible
	Domain() string

	// Explain returns the documentation for a complete identifier e.g.
	// codemark:resource:option. This is used for self-explanatory usage.
	Explain(ident string) string

	// Ressources supported by this generator
	Ressources() []string

	// Generate the artificats based on the given information
	Generate(infos map[*packages.Package]*loaderapi.Information) error

	// OptionsOf returns the options for a choosen resource e.g.
	// codemark:resource.
	OptionsOf(resource string) []core.Option

	// Registry containing all the definitions
	Registry() registry.Registry
}
