package v1

import (
	genv1 "github.com/naivary/codemark/api/generator/v1"
)

type Outputer interface {
	Name() string

	Output(artifacts []*genv1.Artifact, args ...string) error

	Explain() string
}
