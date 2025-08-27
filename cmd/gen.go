package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/outputer"
)

type genCmd struct {
	outputer []string

	outputerMap map[string]string
}

func makeGenCmd(genMngr *generator.Manager, outMngr *outputer.Manager) *cobra.Command {
	g := &genCmd{
		outputerMap: make(map[string]string),
	}
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "generate the artifacts for the given domains",
		Long:    "",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"gen"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			for _, out := range g.outputer {
				domain, outputer, _ := strings.Cut(out, ":")
				g.outputerMap[domain] = outputer
			}
			return g.isValid()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[len(args)-1]
			artifacts, err := genMngr.Generate(pattern)
			if err != nil {
				return err
			}
			for domain, artifacts := range artifacts {
				outputer := g.outputerMap[domain]
				if outputer == "" {
					// TODO: make default outputer configurable
					outputer = "fs"
				}
				err := outMngr.Output(outputer, args, artifacts...)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVarP(&g.outputer, "output", "o", nil, "output destination of the generated artifacts")
	return cmd
}

func (g *genCmd) isValid() error {
	if len(g.outputer) == 0 {
		return errors.New("output cannot be empty")
	}
	return nil
}
