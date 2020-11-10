package server

import (
	"learning-platform/utils"
	"os"
)

type Config struct {
	Domain    string `yaml:"domain"`
	CertsPath string `yaml:"certs_path"`
	CSRF      struct {
		UseCookies    bool     `yaml:"useCookies"`
		ExcludedPaths []string `yaml:"excludedPaths"`
	} `yaml:"csrf"`
}

func (sc *Config) HydrateConfig() {
	sc = utils.HydrateConfig(sc, os.Getenv("CONFIG_PATH")+"/server.yml").(*Config)
}
