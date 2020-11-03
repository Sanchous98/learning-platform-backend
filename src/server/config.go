package server

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type (
	Config interface {
		LoadConfig()
	}

	ServerConfig struct {
		Domain    string `yaml:"domain"`
		CertsPath string `yaml:"certs_path"`
		CSRF      struct {
			UseCookies    bool     `yaml:"useCookies"`
			ExcludedPaths []string `yaml:"excludedPaths"`
		} `yaml:"csrf"`
	}

	GraphQLConfig struct {
		SchemaPath string `yaml:"schema_path"`
	}
)

func (sc *ServerConfig) LoadConfig() {
	sc = loadConfig(sc, os.Getenv("CONFIG_PATH")+"/server.yml").(*ServerConfig)
}

func (gqlc *GraphQLConfig) LoadConfig() {
	gqlc = loadConfig(gqlc, os.Getenv("CONFIG_PATH")+"/graphql.yml").(*GraphQLConfig)
}

func loadConfig(c Config, path string) Config {
	config, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(config, c)

	if err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}

	return c
}
