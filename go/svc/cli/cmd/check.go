package cmd

import (
	"github.com/iv-p/apid/common/config"
	"github.com/iv-p/apid/common/log"
	"github.com/iv-p/apid/svc/cli/http"
	"github.com/iv-p/apid/svc/cli/step"
	"github.com/iv-p/apid/svc/cli/transaction"
	"github.com/iv-p/apid/svc/cli/variables"
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

func checkRun(*cobra.Command, []string) {
	c, err := config.Load(configFilepath)
	if err != nil {
		log.L.Fatalf("could not load config file: %v", err)
	}

	httpClient := http.NewTimedClient()

	stepExecutor := step.NewRequestExecutor(httpClient)
	stepValidator := step.NewResponseValidator()
	stepChecker := step.NewStepChecker(stepExecutor, stepValidator)

	stepInterpolator := transaction.NewStepInterpolator()
	transactionChecker := transaction.NewStepChecker(stepChecker, stepInterpolator)
	transactionService := transaction.NewTransactionService(transactionChecker)

	vars := variables.NewVariables()
	vars = vars.Merge("variables", c.Variables)
	res := transactionService.Check(c.Transactions, vars)
	log.L.Debug(res)
}
