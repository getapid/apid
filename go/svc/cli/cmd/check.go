package cmd

import (
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:     "check",
	Short:   "runs the default config in apid.yaml or the one you provide",
	Long:    "", // TODO maybe this as well
	Example: "", // TODO this
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}
