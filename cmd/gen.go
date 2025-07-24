package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	"github.com/naivary/codemark/generator"
)

func mustInit(fn func() (genv1.Generator, error)) genv1.Generator {
	gen, err := fn()
	if err != nil {
		panic(err)
	}
	return gen
}

type genCmd struct {
	domains []string
}

func makeGenCmd(mngr *generator.Manager) *cobra.Command {
	g := &genCmd{}
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "",
		Long:    "",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"gen"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return g.isValid()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[len(args)-1]
			return mngr.Generate(pattern, g.domains...)
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
