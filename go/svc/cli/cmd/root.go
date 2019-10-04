package cmd

import (
	"fmt"
	"os"

	"github.com/iv-p/apid/common/log"
	"github.com/spf13/cobra"
)

var logVerbosity int

var rootCmd = &cobra.Command{
	Use:  "apid",
	Long: "Apid is a command to help you test and verify the performance of you APIs",
}

func init() {
	rootCmd.PersistentFlags().IntVar(&logVerbosity, "v", -1, "log verbosity")
}

func Execute() {
	log.Init(logVerbosity)
	defer log.L.Sync()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
