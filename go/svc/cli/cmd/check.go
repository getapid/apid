package cmd

import (
	"github.com/iv-p/apid/common/config"
	"github.com/iv-p/apid/common/http"
	"github.com/iv-p/apid/common/log"
	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/common/transaction"
	"github.com/iv-p/apid/common/variables"
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

	err = config.Validate(c)
	if err != nil {
		log.L.Panic("the config failed validation: ", err)
	}

	httpClient := http.NewTimedClient(http.DefaultClient)

	stepInterpolator := step.NewTemplateInterpolator()
	stepExecutor := step.NewHTTPExecutor(httpClient)
	stepValidator := step.NewHTTPValidator()
	stepChecker := step.NewRunner(stepExecutor, stepValidator, stepInterpolator)

	transactionChecker := transaction.NewStepChecker(stepChecker)
	transactionService := transaction.NewTransactionService(transactionChecker)

	vars := variables.NewVariables()
	vars = vars.Merge("variables", c.Variables)
	res := transactionService.Check(c.Transactions, vars)
	log.L.Debug(res)
}
