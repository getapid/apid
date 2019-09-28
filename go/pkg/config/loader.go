package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Loader interface {
	Load() (Config, error)
}

type FileLoader struct {
	Path string
}

func (f FileLoader) Load() (Config, error) {
	cfg := Config{}

	contents, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(contents, &cfg)
	return cfg, err
}
