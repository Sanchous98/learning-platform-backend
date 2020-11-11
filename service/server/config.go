package server

import (
	"github.com/Sanchous98/project-confucius-backend"
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

func (sc *Config) HydrateConfig() error {
	config, err := confucius.HydrateConfig(sc, os.Getenv("CONFIG_PATH")+"/server.yml")

	if err != nil {
		return err
	}

	sc = config.(*Config)

	return nil
}
