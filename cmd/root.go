package cmd

import (
	"fmt"
	"os"

	"github.com/getapid/apid/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:  "apid",
	Long: "Apid is a command to help you test and verify the performance and reliability of you APIs",
	PersistentPreRun: func(*cobra.Command, []string) {
		if verbose {
			log.Atom.SetLevel(zapcore.DebugLevel)
		}
	},
	PersistentPostRun: func(*cobra.Command, []string) {
		log.L.Sync()
	},
}

func init() {
	// default log level for zap is computed in the persistent pre run for the root command
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbosity", "v", false, "output additional logs")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
