package main

import "github.com/iv-p/apid/svc/cli/cmd"

func main() {
<<<<<<< HEAD
	cmd.Execute()
=======
	var configFileLocation = flag.String("c", defaultConfigFileLocation, "location of the config yaml file")
	var defaultLogLevel = flag.Int("v", -1, "default log level")
	flag.Parse()

	log.Init(*defaultLogLevel)
	log.L.Debug("starting apid")
	defer log.L.Sync()

	c, err := config.Load(*configFileLocation)
	if err != nil {
		log.L.Fatalf("could not load config file: %v", err)
	}

	err = config.Validate(c)
	if err != nil {
		log.L.Panic("the config failed validation: ", err)
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
>>>>>>> Config validation impl v1
}
