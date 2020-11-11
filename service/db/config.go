package db

import (
	"github.com/Sanchous98/project-confucius-backend"
	"os"
)

type Config struct {
	DefaultConnection string
	Connections       map[string]interface{}
}

func (db *Config) HydrateConfig() {
	db = confucius.HydrateConfig(db, os.Getenv("CONFIG_PATH")+"/graphql.yml").(*Config)
}
