package v1

import (
	"github.com/spf13/pflag"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
)

type Outputer interface {
	// Documentation for the ouptuter including the name which is unique
	Doc() docv1.Outputer

	Flags() *pflag.FlagSet

	// Output is writing all artificats to the output. `args` are all arguments
	// provided by the user and can be used by the outputer for configuration.
	Output(artifacts []*genv1.Artifact, args ...string) error
}
