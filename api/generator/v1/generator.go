package v1

import (
	"io"

	infov1 "github.com/naivary/codemark/api/info/v1"
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
	Generate(proj infov1.Project) ([]*Artifact, error)

	// Registry containing all the options.
	Registry() registry.Registry
}

type Artifact struct {
	// Name of the Artifact. Make sure the name can be used as a file name
	// because it might be written to the filesystem (including the extension).
	Name string

	Data io.ReadWriter
}
