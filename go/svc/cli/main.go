package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/iv-p/apiping/svc/cli/http"
	"github.com/iv-p/apiping/svc/cli/step"
	"github.com/iv-p/apiping/svc/cli/transaction"
	"github.com/iv-p/apiping/svc/cli/variables"

	"github.com/iv-p/apiping/svc/cli/config"
	"gopkg.in/yaml.v3"
)

func main() {
	var configFileLocation = flag.String("c", config.DefaultConfigFileLocation, "location of the config yaml file")
	flag.Parse()

	var c config.Config
	cfd, err := ioutil.ReadFile(*configFileLocation)
	err = yaml.Unmarshal([]byte(cfd), &c)
	if err != nil {
		log.Fatalf("could not load config file: %v", err)
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
	log.Print(res)
}
