package cloud

import (
	"github.com/spf13/cobra"
)

var suiteCommand = &cobra.Command{
	Use:   "suite",
	Short: "Commands for suite manipulation",
	Long:  `Subcommands include uploading new suite, overwriting existing one, enabling suite, etc.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// noop
	},
}

func init() {
	RootCommand.AddCommand(suiteCommand)
}
