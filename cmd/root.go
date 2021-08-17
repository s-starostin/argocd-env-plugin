package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand returns a new instance of the root command
func NewRootCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "argocd-env-plugin",
		Short: "This is a plugin to add aditional manifests",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	generateCommand := NewGenerateCommand()
	command.AddCommand(generateCommand)
    generateCommand.Flags().StringP("path", "p", "", "a path to directory with manifests")
	generateCommand.Flags().StringP("include-paths", "i", "", "additional paths with manifests")
	command.AddCommand(NewVersionCommand())

	return command
}
