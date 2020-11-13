package graphql

import (
	"github.com/Sanchous98/project-confucius-backend/utils"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	SchemaPath string `yaml:"schema_path"`
}

func (gqlc *Config) HydrateConfig() error {
	config, err := utils.HydrateConfig(gqlc, os.Getenv("CONFIG_PATH")+"/graphql.yml", yaml.Unmarshal)

	if err != nil {
		return err
	}

	gqlc = config.(*Config)

	return nil
}
