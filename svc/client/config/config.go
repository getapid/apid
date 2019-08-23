package config

import (
	"github.com/iv-p/apiping/svc/client/transaction"
	"github.com/iv-p/apiping/svc/client/variables"
)

const (
	// AppName is the name of the project
	AppName = "apid"
	// DefaultConfigFileLocation is the location of the yaml file holding all the config
	DefaultConfigFileLocation = "/etc/" + AppName + "/config.yaml"
)

// Config holds all the config data from the config yaml file
type Config struct {
	APIKey       string                    `yaml:"apikey"`
	Variables    variables.Variables       `yaml:"variables"`
	Transactions []transaction.Transaction `yaml:"transactions"`
}
