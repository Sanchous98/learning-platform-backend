package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config interface {
	HydrateConfig() error
}

func HydrateConfig(c Config, path string) (Config, error) {
	config, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(config, c)

	if err != nil {
		return nil, err
	}

	return c, nil
}
