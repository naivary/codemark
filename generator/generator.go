package generator

import (
	"golang.org/x/tools/go/packages"

	loaderapi "github.com/naivary/codemark/api/loader"
	"github.com/naivary/codemark/registry"
)

type Generator interface {
	// Domain for which the generator is responsible
	Domain() string

	// Explain returns the documentation for an identifier. The identifier can
	// be fully qualified e.g. domain:resource:option or it may be given
	// partially e.g. domain:resource or domain.
	Explain(ident string) string

	// Generate the artificats based on the given information.
	Generate(infos map[*packages.Package]*loaderapi.Information) error

	// Registry containing all the options.
	Registry() registry.Registry
}
