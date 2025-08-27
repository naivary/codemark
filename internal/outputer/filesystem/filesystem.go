package filesystem

import (
	"os"
	"path/filepath"

	"github.com/spf13/pflag"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	outv1 "github.com/naivary/codemark/api/outputer/v1"
)

var _ outv1.Outputer = (*outputer)(nil)

type outputer struct {
	path string
}

func (o *outputer) Name() string {
	return "fs"
}

func (o *outputer) Output(artifacts []*genv1.Artifact, args ...string) error {
	for _, artifact := range artifacts {
		if err := o.output(artifact, args...); err != nil {
			return err
		}
	}
	return nil
}

func (o *outputer) flags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("fs", pflag.ExitOnError)
	flagSet.StringVar(&o.path, "path", "", "path of the location to store the generated artifacts")
	return flagSet
}

func (o *outputer) output(artifact *genv1.Artifact, args ...string) error {
	err := o.flags().Parse(args)
	if err != nil {
		return err
	}
	err = os.MkdirAll(o.path, os.ModeDir)
	if err != nil {
		return err
	}
	filePath := filepath.Join(o.path, artifact.Name)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.ReadFrom(artifact.Data)
	return err
}
