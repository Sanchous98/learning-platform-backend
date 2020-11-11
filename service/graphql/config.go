package graphql

import (
	"github.com/Sanchous98/project-confucius-backend"
	"os"
)

type Config struct {
	SchemaPath string `yaml:"schema_path"`
}

func (gqlc *Config) HydrateConfig() error {
	config, err := confucius.HydrateConfig(gqlc, os.Getenv("CONFIG_PATH")+"/graphql.yml")

	if err != nil {
		return err
	}

	gqlc = config.(*Config)

	return nil
}
