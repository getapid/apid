package config

import (
	"io/ioutil"

	"github.com/iv-p/apid/common/transaction"
	"gopkg.in/yaml.v3"
)

func Load(path string) (Config, error) {
	cfg := config{}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(contents, &cfg)
	if err != nil {
		return Config{}, err
	}

	applySSLFlag(cfg.SkipSSLVerification, cfg.Transactions)

	return Config{
		Version:      cfg.Version,
		APIKey:       cfg.APIKey,
		Variables:    cfg.Variables,
		Transactions: cfg.Transactions,
	}, err
}

func applySSLFlag(skipSSL bool, txs []transaction.Transaction) {
	for _, tx := range txs {
		for j := range tx.Steps {
			tx.Steps[j].Request.SkipSSLVerification = skipSSL
		}
	}
}
