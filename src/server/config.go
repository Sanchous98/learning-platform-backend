package server

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config interface {
	LoadConfig()
}

type ServerConfig struct {
	Domain    string `yaml:"domain"`
	CertsPath string `yaml:"certs_path"`
	CSRF      struct {
		UseCookies    bool     `yaml:"useCookies"`
		ExcludedPaths []string `yaml:"excludedPaths"`
	} `yaml:"csrf"`
}

func (sc *ServerConfig) LoadConfig() {
	config, _ := ioutil.ReadFile(os.Getenv("CONFIG_PATH") + "/server.yml")
	err := yaml.Unmarshal(config, sc)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
