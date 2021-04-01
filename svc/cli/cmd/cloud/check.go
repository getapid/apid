package cloud

import (
	"fmt"
	"os"

	native "net/http"

	"github.com/getapid/apid-cli/common/config"
	"github.com/getapid/apid-cli/common/log"
	"github.com/getapid/apid-cli/common/result"
	"github.com/getapid/apid-cli/common/step"
	"github.com/getapid/apid-cli/common/transaction"
	"github.com/getapid/apid-cli/common/variables"
	cmdResult "github.com/getapid/apid-cli/svc/cli/result"
	"github.com/spf13/cobra"
)

var (
	region      string
	showTimings bool
	parallelism int
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Executes a config on apid remote infrastructure",
	Long: `Loads either a config file, or if a directory is provided it 
recursively loads all .yaml files and executes them. It will run all the 
transactions in you config, verify the responses, and record the time it 
took to action each request. Runs all steps in each transaction on a 
public cloud infrastructure. See https://docs.getapid.com for more inforamtion`,
	Example: `
	apid cloud check --key <apid access key>
	apid cloud check --config my-api.yaml --key <apid access key> --region dublin
	apid cloud check -c ./e2e-tests/ -k <apid access key>
	apid cloud check -c ./e2e-tests/ -k <apid access key> -r dublin`,
	Args: cobra.NoArgs,
	RunE: remoteCheck,
}

func init() {
	RootCommand.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&configFilepath, "config", "c", "./apid.yaml", "file with config to run")
	checkCmd.Flags().BoolVarP(&showTimings, "timings", "t", false, "output the durations of requests")
	checkCmd.Flags().IntVarP(&parallelism, "parallelism", "p", 10, "number of concurrent transaction executions")
	checkCmd.Flags().StringVarP(&region, "region", "r", "washington", "location to run the tests from")
}

func remoteCheck(cmd *cobra.Command, args []string) error {
	c, err := config.Load(configFilepath)
	if err != nil {
		return fmt.Errorf("could not load config file: %v", err)
	}

	err = config.Validate(c)
	if err != nil {
		return fmt.Errorf("the config failed validation: %v", err)
	}

	consoleWriter := cmdResult.NewConsoleWriter(cmd.OutOrStdout(), showTimings)
	writer := result.NewMultiWriter(consoleWriter)

	stepInterpolator := step.NewTemplateInterpolator()
	stepExecutor, err := step.NewRemoteHTTPExecutor(native.DefaultClient, apiKey, region)
	if err != nil {
		log.L.Error(err)
		os.Exit(1)
	}
	stepValidator := step.NewHTTPValidator()
	stepExtractor := step.NewBodyExtractor()
	stepChecker := step.NewRunner(stepExecutor, stepValidator, stepInterpolator, stepExtractor)

	transactionRunner := transaction.NewTransactionRunner(stepChecker, writer, parallelism)

	vars := c.Variables.Merge(variables.New(variables.WithEnv()))
	ok := transactionRunner.Run(c.Transactions, vars)
	if !ok {
		os.Exit(1)
	}
	return nil
}
