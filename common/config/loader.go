package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/getapid/cli/common/transaction"
	"gopkg.in/yaml.v3"
)

func LoadReader(r io.Reader) (Config, error) {
	return loadReader(r)
}

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
	suiteFile, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer suiteFile.Close()

	return loadReader(suiteFile)
}

func loadReader(r io.Reader) (Config, error) {
	cfg := &config{}
	err := yaml.NewDecoder(r).Decode(cfg)
	if err != nil {
		if err == io.EOF {
			return Config{}, nil
		}
		return Config{}, err
	}

	applySSLFlag(cfg.SkipSSLVerification, cfg.Transactions)

	return Config{
		Version:      cfg.Version,
		APIKey:       cfg.APIKey,
		Variables:    cfg.Variables,
		Schedule:     cfg.Schedule,
		Locations:    cfg.Locations,
		Transactions: cfg.Transactions,
	}, nil
}

func mergeConfigs(configs []Config) (Config, error) {
	result := Config{}
	var err error = nil

	for _, other := range configs {
		result.Variables = result.Variables.Merge(other.Variables)
		result.Transactions = append(result.Transactions, other.Transactions...)
		if other.Schedule != "" {
			result.Schedule = other.Schedule
		}
		if len(other.Locations) != 0 {
			result.Locations = other.Locations
		}
	}
	return result, err
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
			if tx.Steps[j].Request.SkipSSLVerification == nil {
				tx.Steps[j].Request.SkipSSLVerification = &skipSSL
			}
		}
	}
}
