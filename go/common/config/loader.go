package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/iv-p/apid/common/transaction"
	"gopkg.in/yaml.v3"
)

func Load(path string) (Config, error) {
	return tryLoad(path)
}

func tryLoad(path string) (Config, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return Config{}, fmt.Errorf("loading %s: %w", path, err)
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return loadDir(path)
	case mode.IsRegular():
		return loadFile(path)
	}
	return Config{}, nil
}

func loadDir(path string) (Config, error) {
	var files []string
	err := filepath.Walk(path, func(p string, i os.FileInfo, err error) error {
		if p != path && isFile(p) && (strings.HasSuffix(p, ".yaml") || strings.HasSuffix(p, ".yml")) {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return Config{}, err
	}

	configs := make([]Config, len(files))
	for i, file := range files {
		config, err := loadFile(file)
		if err != nil {
			return Config{}, err
		}
		configs[i] = config
	}
	return mergeConfigs(configs)
}

func loadFile(path string) (Config, error) {
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

func mergeConfigs(configs []Config) (Config, error) {
	result := Config{}
	var err error = nil

	for _, other := range configs {
		result.Variables = mergeVariables(result.Variables, other.Variables)
		result.Transactions = append(result.Transactions, other.Transactions...)
	}
	return result, err
}

func mergeVariables(this, other map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range this {
		result[key] = value
	}
	for key, value := range other {
		result[key] = value
	}
	return result
}

func isFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		return true
	}
	return false
}

func applySSLFlag(skipSSL bool, txs []transaction.Transaction) {
	for _, tx := range txs {
		for j := range tx.Steps {
			tx.Steps[j].Request.SkipSSLVerification = skipSSL
		}
	}
}
