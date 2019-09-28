package main

import (
	"flag"
	"io/ioutil"

	"github.com/iv-p/apid/pkg/log"
	"github.com/iv-p/apid/svc/cli/config"
	"github.com/iv-p/apid/svc/cli/http"
	"github.com/iv-p/apid/svc/cli/step"
	"github.com/iv-p/apid/svc/cli/transaction"
	"github.com/iv-p/apid/svc/cli/variables"
	"gopkg.in/yaml.v3"
)

func main() {
	var configFileLocation = flag.String("c", config.DefaultConfigFileLocation, "location of the config yaml file")
	var defaultLogLevel = flag.Int("v", -1, "default log level")
	flag.Parse()

	log.Init(*defaultLogLevel)
	log.L.Debug("starting apid")
	defer log.L.Sync()

	var c config.Config
	cfd, err := ioutil.ReadFile(*configFileLocation)
	err = yaml.Unmarshal([]byte(cfd), &c)
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
