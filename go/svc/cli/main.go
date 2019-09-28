package main

import (
	"flag"

	"github.com/iv-p/apid/pkg/config"
	"github.com/iv-p/apid/pkg/log"
	"github.com/iv-p/apid/svc/cli/http"
	"github.com/iv-p/apid/svc/cli/step"
	"github.com/iv-p/apid/svc/cli/transaction"
	"github.com/iv-p/apid/svc/cli/variables"
)

// defaultConfigFileLocation is the location of the yaml file holding all the config
const defaultConfigFileLocation = "config.yaml"

func main() {
	var configFileLocation = flag.String("c", defaultConfigFileLocation, "location of the config yaml file")
	var defaultLogLevel = flag.Int("v", -1, "default log level")
	flag.Parse()

	log.Init(*defaultLogLevel)
	log.L.Debug("starting apid")
	defer log.L.Sync()

	loader := config.FileLoader{Path: *configFileLocation}
	c, err := loader.Load()
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
