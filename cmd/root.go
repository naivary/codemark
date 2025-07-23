package cmd

import (
	"slices"

	"github.com/spf13/cobra"

	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/internal/generator/k8s"
)

// TODO: make the fu nctinos of the commands use codes
const (
	Success = iota
	InternalErr
	BadRequest
)

func Exec(convs []converter.Converter, gens []generator.Generator) (int, error) {
	mngr, err := newGenManager(gens)
	if err != nil {
		return InternalErr, err
	}
	rootCmd := makeRootCmd()
	rootCmd.AddCommand(
		makeGenCmd(mngr),
		makeExplainCmd(mngr),
	)
	err = rootCmd.Execute()
	if err != nil {
		return 1, err
	}
	return Success, nil
}

func makeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "codemark",
		Short:        "codemark is a annotation library for go allowing you to generate any kind of artifats.",
		Long:         ``,
		SilenceUsage: true,
	}
	return cmd
}

func newGenManager(gens []generator.Generator) (*generator.Manager, error) {
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
