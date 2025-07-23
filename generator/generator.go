package generator

import (
	"golang.org/x/tools/go/packages"

	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/registry"
)

type Generator interface {
	// Domain for which the generator is responsible
	Domain() string

	// Explain returns the documentation for a complete identifier e.g.
	// codemark:resource:option. This is used for self-explanatory usage.
	Explain(ident string) string

	// Generate the artificats based on the given information
	Generate(infos map[*packages.Package]*loaderapi.Information) error

	// Registry containing all the definitions
	Registry() registry.Registry
}
