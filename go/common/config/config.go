package config

import (
	"github.com/iv-p/apid/common/transaction"
)

// Config holds the internal representation of the config
type Config struct {
	Version      string
	APIKey       string
	Variables    map[string]interface{}
	Transactions []transaction.Transaction
}

// Config holds all the config data from the config yaml file
type config struct {
	Version             string                    `yaml:"version" validate:"version"`
	APIKey              string                    `yaml:"apikey"`
	Variables           map[string]interface{}    `yaml:"variables"`
	Transactions        []transaction.Transaction `yaml:"transactions" validate:"required,unique=ID"`
	SkipSSLVerification bool                      `yaml:"skip_ssl_verify"`
}
