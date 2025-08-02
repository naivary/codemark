package cmd

import (
	"slices"

	"github.com/spf13/cobra"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	"github.com/naivary/codemark/generator"
	k8sgen "github.com/naivary/codemark/internal/generator/k8s"
)

// TODO: make the functions of the commands use codes
const (
	Success = iota
	InternalErr
	BadRequest
)

func Exec(convs []convv1.Converter, gens []genv1.Generator) (int, error) {
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

func mustInit(fn func() (genv1.Generator, error)) genv1.Generator {
	gen, err := fn()
	if err != nil {
		panic(err)
	}
	return gen
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

func newGenManager(gens []genv1.Generator) (*generator.Manager, error) {
	mngr, err := generator.NewManager()
	if err != nil {
		return nil, err
	}
	builtinGens := []genv1.Generator{
		mustInit(k8sgen.New),
	}
	for _, gen := range slices.Concat(builtinGens, gens) {
		if err := mngr.Add(gen); err != nil {
			return nil, err
		}
	}
	return mngr, nil
}
