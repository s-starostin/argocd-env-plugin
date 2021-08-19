package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand returns a new instance of the root command
func NewRootCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "argocd-env-plugin",
		Short: "Merge manifests from different paths",
		Long: "This is an Argo CD plugin to merge multiple manifests from different paths together.\r\n" +
			"Instead of specifing CLI params for primary and additional manifest paths," +
			"you can use environment variables AEP_PATH and AEP_INCLUDE_PATHS.",
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
