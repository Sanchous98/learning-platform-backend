package server

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
}

func GetConfig() {
	err := yaml.Unmarshal()
}
