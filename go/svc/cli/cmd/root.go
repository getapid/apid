package cmd

import (
	"fmt"
	"os"

	"github.com/getapid/apid/common/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logVerbosity int

var rootCmd = &cobra.Command{
	Use:  "apid",
	Long: "Apid is a command to help you test and verify the performance of you APIs",
	PersistentPreRun: func(*cobra.Command, []string) {
		log.Init(int(zap.ErrorLevel) - logVerbosity)
	},
	PersistentPostRun: func(*cobra.Command, []string) {
		log.L.Sync()
	},
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&logVerbosity, "verbosity", "v", 0, "log verbosity (default: Error)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
