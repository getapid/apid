package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func Load(path string) (Config, error) {
	cfg := Config{}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(contents, &cfg)
	return cfg, err
}
