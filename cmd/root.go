package main

import "github.com/spf13/cobra"

func makeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "codemark",
		Short: "codemark is a annotation library for go allowing you to generate any kind of artifats.",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return cmd
}
