package outputer

import (
	"io"
	"os"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	outv1 "github.com/naivary/codemark/api/outputer/v1"
)

var _ outv1.Outputer = (*fsOutputer)(nil)

type stdoutOutputer struct{}

func NewStdoutOutputer() (outv1.Outputer, error) {
	return &stdoutOutputer{}, nil
}

func (o *stdoutOutputer) Name() string {
	return "stdout"
}

func (o *stdoutOutputer) Output(artifacts []*genv1.Artifact, args ...string) error {
	for _, artifact := range artifacts {
		if err := o.output(artifact); err != nil {
			return err
		}
	}
	return nil
}

func (o *stdoutOutputer) output(artifact *genv1.Artifact) error {
	_, err := io.Copy(os.Stdout, artifact.Data)
	return err
}
