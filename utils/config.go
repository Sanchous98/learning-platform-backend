package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config is basic interface for every service configuration
type Config interface {
	HydrateConfig() error
}

// HydrateConfig is a basic function to read configuration from a yaml
func HydrateConfig(c Config, path string) (Config, error) {
	config, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(config, c)

	if err != nil {
		return nil, err
	}

	return c, nil
}
