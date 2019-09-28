package config

import (
	"github.com/iv-p/apid/svc/cli/transaction"
)

const (
	// AppName is the name of the project
	AppName = "apid"
)

// Config holds all the config data from the config yaml file
type Config struct {
	APIKey       string                    `yaml:"apikey"`
	Variables    map[string]interface{}    `yaml:"variables"`
	Transactions []transaction.Transaction `yaml:"transactions"`
}
