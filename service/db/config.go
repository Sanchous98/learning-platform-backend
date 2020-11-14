package db

import (
	"github.com/Sanchous98/project-confucius-backend/utils"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	DefaultConnection string
	Connections       map[string]interface{}
}

func (db *Config) HydrateConfig() error {
	config, err := utils.HydrateConfig(db, os.Getenv("CONFIG_PATH")+"/database.yml", yaml.Unmarshal)

	if err != nil {
		return err
	}

	db = config.(*Config)

	return nil
}
