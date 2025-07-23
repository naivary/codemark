package cmd

import (
	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/generator"
)

func Exec(convs []converter.Converter, gens []generator.Generator) (int, error) {
	rootCmd := makeRootCmd()
	rootCmd.AddCommand(
		makeGenCmd(gens...),
	)
	err := rootCmd.Execute()
	if err != nil {
		return 1, err
	}
	return 0, nil
}
