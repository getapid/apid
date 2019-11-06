package cmd

import (
	"os"

	"github.com/iv-p/apid/common/config"
	"github.com/iv-p/apid/common/http"
	"github.com/iv-p/apid/common/log"
	"github.com/iv-p/apid/common/result"
	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/common/transaction"
	"github.com/iv-p/apid/common/variables"
	cmdResult "github.com/iv-p/apid/svc/cli/result"
	"github.com/spf13/cobra"
)

var configFilepath string

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Executes a config",
	Long: `Check gets the default config or the one you have provided and executes it.
It will run all the transactions in you config, verify the responses,
and record the time it took to action each request`,
	Example: `
	apid check	
	apid check --config my-api.yaml`,
	Args: cobra.NoArgs,
	Run:  checkRun,
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&configFilepath, "config", "c", "./apid.yaml", "file with config to run")
}

func checkRun(cmd *cobra.Command, args []string) {
	c, err := config.Load(configFilepath)
	if err != nil {
		log.L.Fatalf("could not load config file: %v", err)
	}

	err = config.Validate(c)
	if err != nil {
		log.L.Panic("the config failed validation: ", err)
	}

	consoleWriter := cmdResult.NewConsoleWriter(cmd.OutOrStdout())
	writer := result.NewMultiWriter(consoleWriter)

	httpClient := http.NewTimedClient(http.DefaultClient)

	stepInterpolator := step.NewTemplateInterpolator()
	stepExecutor := step.NewHTTPExecutor(httpClient)
	stepValidator := step.NewHTTPValidator()
	stepExtractor := step.NewBodyExtractor()
	stepChecker := step.NewRunner(stepExecutor, stepValidator, stepInterpolator, stepExtractor)

	transactionRunner := transaction.NewTransactionRunner(stepChecker, writer)

	vars := variables.New(variables.WithVars(c.Variables), variables.WithEnv())
	ok := transactionRunner.Run(c.Transactions, vars)
	if !ok {
		os.Exit(1)
	}
}
