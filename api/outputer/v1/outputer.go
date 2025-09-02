package v1

import (
	genv1 "github.com/naivary/codemark/api/generator/v1"
)

type Outputer interface {
	// Name of the outputer. This name should be unique within the project
	// because it is used for indexing.
	Name() string

	// Output is writing all artificats to the output. `args` are all arguments
	// provided by the user and can be used by the outputer for configuration.
	Output(artifacts []*genv1.Artifact, args ...string) error
}
