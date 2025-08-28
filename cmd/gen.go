package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/outputer"
)

type genCmd struct {
	outputers []string
}

func makeGenCmd(cfg *cliConfig, genMngr *generator.Manager, outMngr *outputer.Manager, convs []convv1.Converter) *cobra.Command {
	g := &genCmd{}
	cmd := &cobra.Command{
		Use:     "generate [pattern]",
		Short:   "generate the artifacts for the given pattern",
		Aliases: []string{"gen"},
		RunE:    g.runE(cfg, genMngr, outMngr, convs),
	}
	return cmd
}

func (g *genCmd) runE(
	cfg *cliConfig,
	genMngr *generator.Manager,
	outMngr *outputer.Manager,
	convs []convv1.Converter,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		pattern := args[0]
		artifacts, err := genMngr.Generate(convs, pattern)
		if err != nil {
			return err
		}
		outMap := g.outputerMap(cfg, genMngr.Domains())
		for domain, artifacts := range artifacts {
			err := outMngr.Output(outMap[domain], args[1:], artifacts...)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (g *genCmd) outputerMap(cfg *cliConfig, domains []string) map[string]string {
	const sep = ":"
	res := make(map[string]string, len(g.outputers))
	for _, domain := range domains {
		res[domain] = cfg.DefaultOutputer
	}
	for _, outputer := range g.outputers {
		domain, outputerName, found := strings.Cut(outputer, sep)
		if !found {
			continue
		}
		res[domain] = outputerName
	}
	return res
}
