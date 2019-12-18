package config

import (
	"github.com/getapid/apid-cli/common/transaction"
	"github.com/getapid/apid-cli/common/variables"
)

// Config holds the internal representation of the config
type Config struct {
	Version      string
	APIKey       string
	Variables    variables.Variables
	Transactions []transaction.Transaction
}

// Config holds all the config data from the config yaml file
type config struct {
	Version             string                    `yaml:"version"`
	APIKey              string                    `yaml:"apikey"`
	Variables           variables.Variables       `yaml:"variables"`
	Transactions        []transaction.Transaction `yaml:"transactions" validate:"unique=ID"`
	SkipSSLVerification bool                      `yaml:"skip_ssl_verify"`
}
