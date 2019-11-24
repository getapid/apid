package cmd

import (
	"fmt"
	"os"

	"github.com/getapid/apid/common/config"
	"github.com/getapid/apid/common/http"
	"github.com/getapid/apid/common/result"
	"github.com/getapid/apid/common/step"
	"github.com/getapid/apid/common/transaction"
	"github.com/getapid/apid/common/variables"
	cmdResult "github.com/getapid/apid/svc/cli/result"
	"github.com/spf13/cobra"
)

var configFilepath string

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Executes a config",
	Long: `Loads either a config file, or if a directory is provided it 
recursively loads all .yaml files and executes them. It will run all the 
transactions in you config, verify the responses, and record the time it 
took to action each request`,
	Example: `
	apid check	
	apid check --config my-api.yaml
	apid check --config ./e2e-tests/`,
	Args: cobra.NoArgs,
	RunE: checkRun,
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&configFilepath, "config", "c", "./apid.yaml", "file with config to run")
}

func checkRun(cmd *cobra.Command, args []string) error {
	c, err := config.Load(configFilepath)
	if err != nil {
		return fmt.Errorf("could not load config file: %v", err)
	}

	err = config.Validate(c)
	if err != nil {
		return fmt.Errorf("the config failed validation: %v", err)
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
	return nil
}
