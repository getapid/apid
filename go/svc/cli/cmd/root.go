package cmd

import (
	"fmt"
	"os"

	"github.com/iv-p/apid/common/log"
	"github.com/spf13/cobra"
)

var (
	configFilepath string
	logVerbosity   int

	RootCmd = &cobra.Command{
		Use:   "apid",
		Short: "Apid is a command to help you test and verify the performance of you API",
		Long: `The tool will make sure all your API's endpoints agree with the spec you have provided.
Along with that it will record the time it took for your app to complete the requests and break down its
lifecycle so you know where to improve.`,
		PersistentPostRun: func(*cobra.Command, []string) {
			_ = log.L.Sync()
		},
	}
)

func init() {
	cobra.OnInitialize(setupLog)

	RootCmd.AddCommand(checkCmd)

	RootCmd.PersistentFlags().StringVar(&configFilepath, "config", "./apid.yaml", "file with config to run")
	RootCmd.PersistentFlags().IntVar(&logVerbosity, "v", -1, "log verbosity")
}

func setupLog() {
	log.Init(logVerbosity)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
