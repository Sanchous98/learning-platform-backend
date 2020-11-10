package authentication

import (
	confucius "github.com/Sanchous98/project-confucius-backend"
	"github.com/Sanchous98/project-confucius-backend/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type authentication struct {
	config *Config
}

func NewService(config utils.Config) confucius.Service {
	return &authentication{config.(*Config)}
}

func (a *authentication) Serve(router *mux.Router) error {

	return nil
}

func (a *authentication) Stop() {

}

func (a *authentication) Init() error {
	return a.config.HydrateConfig()
}

func (a *authentication) Middleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
