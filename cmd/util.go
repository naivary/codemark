package cmd

import (
	"slices"

	"github.com/spf13/cobra"

	genv1 "github.com/naivary/codemark/api/generator/v1"
	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/internal/generator/openapi"
)

type run = func(cmd *cobra.Command, args []string) error

func mustInit(fn func() (genv1.Generator, error)) genv1.Generator {
	gen, err := fn()
	if err != nil {
		panic(err)
	}
	return gen
}

func newGenManager(configPath string, gens []genv1.Generator) (*generator.Manager, error) {
	mngr, err := generator.NewManager(configPath)
	if err != nil {
		return nil, err
	}
	builtinGens := []genv1.Generator{
		mustInit(openapi.New),
	}
	for _, gen := range slices.Concat(builtinGens, gens) {
		if err := mngr.Add(gen); err != nil {
			return nil, err
		}
	}
	return mngr, nil
}
