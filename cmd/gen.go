package cmd

import (
	"errors"
	"slices"

	"github.com/spf13/cobra"

	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/internal/generator/k8s"
)

func mustInit(fn func() (generator.Generator, error)) generator.Generator {
	gen, err := fn()
	if err != nil {
		panic(err)
	}
	return gen
}

type genCmd struct {
	domains []string
}

func makeGenCmd(gens ...generator.Generator) *cobra.Command {
	g := &genCmd{}
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "",
		Long:    "",
		Aliases: []string{"gen"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := g.isValid(); err != nil {
				return err
			}
			gen := g.generate(gens...)
			return gen(cmd, args)
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

func (g *genCmd) generate(gens ...generator.Generator) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		mngr, err := g.newManager(gens...)
		if err != nil {
			return err
		}
		pattern := args[len(args)-1]
		return mngr.Generate(pattern, g.domains...)
	}
}

func (g *genCmd) newManager(gens ...generator.Generator) (*generator.Manager, error) {
	mngr, err := generator.NewManager()
	if err != nil {
		return nil, err
	}
	builtinGens := []generator.Generator{
		mustInit(k8s.NewGenerator),
	}
	for _, gen := range slices.Concat(builtinGens, gens) {
		if err := mngr.Add(gen); err != nil {
			return nil, err
		}
	}
	return mngr, nil
}
