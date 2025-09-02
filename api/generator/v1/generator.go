package v1

import (
	"io"

	infov1 "github.com/naivary/codemark/api/info/v1"
	regv1 "github.com/naivary/codemark/api/registry/v1"
)

type Generator interface {
	// Domain for which the generator is responsible
	Domain() string

	// Generate the artificats based on the given information.
	Generate(proj infov1.Project, config map[string]any) ([]*Artifact, error)

	// Registry containing all the options.
	Registry() regv1.Registry
}

type Artifact struct {
	// Name of the Artifact. Make sure the name can be used as a file name
	// because it might be written to the filesystem (including the extension).
	Name string

	Data io.ReadWriter
}
