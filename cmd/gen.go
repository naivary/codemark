package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/naivary/codemark/generator"
)

type genCmd struct {
	domains []string
	output  string
}

func makeGenCmd(mngr *generator.Manager) *cobra.Command {
	g := &genCmd{}
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "generate the artifacts for the given domains",
		Long:    "",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"gen"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return g.isValid()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[len(args)-1]
			artifacts, err := mngr.Generate(pattern, g.domains...)
			if err != nil {
				return err
			}
			err = os.MkdirAll(g.output, os.ModePerm)
			if err != nil {
				return err
			}
			for _, artifacts := range artifacts {
				for _, artifact := range artifacts {
					path := filepath.Join(g.output, artifact.Name)
					file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o777)
					defer file.Close()
					if err != nil {
						return err
					}
					_, err = file.ReadFrom(artifact.Data)
					if err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&g.domains, "domains", nil, "domains to generate artifacts for")
	cmd.Flags().StringVar(&g.output, "output", "", "output destination of the generated artifacts")
	return cmd
}

func (g *genCmd) isValid() error {
	if len(g.domains) == 0 {
		return errors.New("domains cannot be empty")
	}
	if len(g.output) == 0 {
		return errors.New("output cannot be empty")
	}
	return nil
}
