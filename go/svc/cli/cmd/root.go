package cmd

import (
	"fmt"
	"os"

	"github.com/iv-p/apid/common/log"
	"github.com/spf13/cobra"
)

var logVerbosity int

var RootCmd = &cobra.Command{
	Use:  "apid",
	Long: "Apid is a command to help you test and verify the performance of you APIs",
	PersistentPostRun: func(*cobra.Command, []string) {
		_ = log.L.Sync()
	},
}

func init() {
	cobra.OnInitialize(func() {
		log.Init(logVerbosity) // the value isn't available until after the init is finished
	})

	RootCmd.PersistentFlags().IntVar(&logVerbosity, "v", -1, "log verbosity")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
