package outputer

import (
	"os"
	"path/filepath"

	"github.com/spf13/pflag"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	outv1 "github.com/naivary/codemark/api/outputer/v1"
)

var _ outv1.Outputer = (*fsOutputer)(nil)

type fsOutputer struct {
	path string
}

func NewFsOutputer() (outv1.Outputer, error) {
	return &fsOutputer{}, nil
}

func (o *fsOutputer) Doc() docv1.Outputer {
	return docv1.Outputer{
		Name: "fs",
		Desc: "",
	}
}

func (o *fsOutputer) Output(artifacts []*genv1.Artifact, args ...string) error {
	for _, artifact := range artifacts {
		if err := o.output(artifact, args...); err != nil {
			return err
		}
	}
	return nil
}

func (o *fsOutputer) Flags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("fs", pflag.ContinueOnError)
	flagSet.ParseErrorsWhitelist.UnknownFlags = true
	flagSet.StringVar(&o.path, "fs.path", "", "path of the location to store the generated artifacts")
	return flagSet
}

func (o *fsOutputer) output(artifact *genv1.Artifact, args ...string) error {
	err := o.Flags().Parse(args)
	if err != nil {
		return err
	}
	if o.path == "" {
		o.path, _ = os.Getwd()
		o.path = filepath.Join(o.path, "codemark")
	}
	err = os.MkdirAll(o.path, os.ModePerm)
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
