package v1

import (
	"io"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	infov1 "github.com/naivary/codemark/api/info/v1"
	regv1 "github.com/naivary/codemark/api/registry/v1"
)

type Generator interface {
	// Domain for which the generator is responsible
	Domain() docv1.Domain

	// Generate the artificats based on the given information.
	Generate(proj infov1.Project, config map[string]any) ([]*Artifact, error)

	// Registry containing all the options.
	Registry() regv1.Registry

	// Resources which are supported by the generator with appropiate
	// documentation. The map is indexed by the resource name.
	Resources() map[string]*docv1.Resource

	// ConfigDoc is returning the documentation of the available configuration
	// options. It is used by the explain command to provide self documentation.
	ConfigDoc() map[string]docv1.Config
}

type Artifact struct {
	// Name of the Artifact. Make sure the name can be used as a file name
	// because it might be written to the filesystem (including the extension).
	Name string

	// The actual data of the artifact created by interpreting the markers.
	Data io.ReadWriter
}
