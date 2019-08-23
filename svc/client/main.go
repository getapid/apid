package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/iv-p/apiping/svc/client/http"
	"github.com/iv-p/apiping/svc/client/step"
	"github.com/iv-p/apiping/svc/client/transaction"

	"github.com/iv-p/apiping/svc/client/config"
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

	res := transactionService.Check(c.Transactions, c.Variables)
	log.Print(res)
}
