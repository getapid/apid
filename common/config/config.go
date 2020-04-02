package config

import (
	"github.com/getapid/apid-cli/common/transaction"
	"github.com/getapid/apid-cli/common/variables"
)

// Config holds the internal representation of the config. It has yaml tags to allow to marshalling
type Config struct {
	Version      string                    `yaml:"version"`
	APIKey       string                    `yaml:"api_key"`
	Schedule     string                    `yaml:"schedule" validate:"cron"`
	Locations    []string                  `yaml:"locations"`
	Variables    variables.Variables       `yaml:"variables"`
	Transactions []transaction.Transaction `yaml:"transactions" validate:"unique=ID"`
}

// Config holds all the config data from the config yaml file
type config struct {
	Version             string                    `yaml:"version"`
	APIKey              string                    `yaml:"apikey"`
	Schedule            string                    `yaml:"schedule"`
	Locations           []string                  `yaml:"locations"`
	Variables           variables.Variables       `yaml:"variables"`
	Transactions        []transaction.Transaction `yaml:"transactions"`
	SkipSSLVerification bool                      `yaml:"skip_ssl_verify"`
}
