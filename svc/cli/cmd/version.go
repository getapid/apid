package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of APId",
	Long:  `Print the version number of APId`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("APId %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
