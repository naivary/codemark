package cmd

import (
	"os"

	"github.com/spf13/cobra"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	outv1 "github.com/naivary/codemark/api/outputer/v1"
	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/outputer"
)

// TODO: make the functions of the commands use codes
const (
	Success = iota
	InternalErr
	BadRequest
)

func Exec(gens []genv1.Generator, outs []outv1.Outputer, convs []convv1.Converter) (int, error) {
	rootCmd := makeRootCmd()
	err := rootCmd.ParseFlags(os.Args)
	if err != nil {
		return InternalErr, err
	}
	cfgFile, err := rootCmd.Flags().GetString("config")
	if err != nil {
		return InternalErr, err
	}
	cfg, err := newConfig(cfgFile)
	if err != nil {
		return InternalErr, err
	}
	genMngr, outMngr, err := makeManager(cfgFile, gens, outs)
	if err != nil {
		return InternalErr, err
	}
	rootCmd.AddCommand(
		makeGenCmd(cfg, genMngr, outMngr, convs),
	)
	err = rootCmd.Execute()
	if err != nil {
		return 1, err
	}
	return Success, nil
}

type rootCmd struct {
	cfgFile string
}

func makeRootCmd() *cobra.Command {
	r := &rootCmd{}
	cmd := &cobra.Command{
		Use:          "codemark",
		Short:        "annotation library allowing you to generate anything.",
		SilenceUsage: true,
	}
	cmd.PersistentFlags().StringVar(&r.cfgFile, "config", "", "path of codemark.yaml config file to use")
	return cmd
}

func makeManager(
	cfgFile string,
	gens []genv1.Generator,
	outs []outv1.Outputer,
) (*generator.Manager, *outputer.Manager, error) {
	genMngr, err := newGenManager(cfgFile, gens)
	if err != nil {
		return nil, nil, err
	}
	outMngr, err := newOutManager(outs)
	return genMngr, outMngr, err
}
