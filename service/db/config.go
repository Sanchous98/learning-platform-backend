package db

import (
	"learning-platform/utils"
	"os"
)

type Config struct {
	DefaultConnection string
	Connections       map[string]interface{}
}

func (db *Config) HydrateConfig() {
	db = utils.HydrateConfig(db, os.Getenv("CONFIG_PATH")+"/graphql.yml").(*Config)
}
