package graphql

import (
	"learning-platform/utils"
	"os"
)

type Config struct {
	SchemaPath string `yaml:"schema_path"`
}

func (gqlc *Config) HydrateConfig() {
	gqlc = utils.HydrateConfig(gqlc, os.Getenv("CONFIG_PATH")+"/graphql.yml").(*Config)
}
