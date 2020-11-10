package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config interface {
	HydrateConfig()
}

func HydrateConfig(c Config, path string) Config {
	config, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(config, c)

	if err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}

	return c
}
