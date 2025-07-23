package cmd

import "github.com/spf13/cobra"

func makeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "codemark",
		Short: "codemark is a annotation library for go allowing you to generate any kind of artifats.",
		Long:  ``,
	}
	return cmd
}
