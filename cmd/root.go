package main

import "github.com/spf13/cobra"

func makeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "codemark",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return cmd
}
