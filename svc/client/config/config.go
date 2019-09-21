package config

import (
	"github.com/iv-p/apiping/svc/client/transaction"
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
	Variables    map[string]interface{}    `yaml:"variables"`
	Transactions []transaction.Transaction `yaml:"transactions"`
}
