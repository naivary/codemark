package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/naivary/codemark/generator"
)

type genCmd struct {
	domains []string
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
			_ = artifacts
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&g.domains, "domains", nil, "domains to generate artifacts for")
	return cmd
}

func (g *genCmd) isValid() error {
	if len(g.domains) == 0 {
		return errors.New("domains cannot be empty")
	}
	return nil
}
