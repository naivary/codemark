package main

import (
	"github.com/spf13/cobra"

	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/internal/generator/k8s"
	"github.com/naivary/codemark/loader"
	"github.com/naivary/codemark/registry"
)

type genCmd struct {
	domains []string
}

func makeGenCmd() *cobra.Command {
	g := &genCmd{}
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "",
		Long:    "",
		Aliases: []string{"gen"},
		RunE:    g.generate,
	}
	cmd.Flags().StringSliceVar(&g.domains, "domain", nil, "domains to generate artifacts for")
	return cmd
}

func (g *genCmd) generate(cmd *cobra.Command, args []string) error {
	mngr, err := generator.NewManager()
	if err != nil {
		return err
	}
	gens := []generator.Generator{}
	k8sGen, err := k8s.NewGenerator()
	if err != nil {
		return err
	}
	gens = append(gens, k8sGen)
	if err := mngr.Add(k8sGen); err != nil {
		return err
	}
	reg, err := registry.Merge(k8sGen.Registry())
	if err != nil {
		return err
	}
	pattern := args[len(args)-1]
	info, err := loader.Load(reg, pattern)
	if err != nil {
		return err
	}
	// merge registries
	// load infos
	// call generate for every generator
	for _, gen := range gens {
		if err := gen.Generate(info); err != nil {
			return err
		}
	}
	return nil
}
