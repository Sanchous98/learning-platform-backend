package server

import (
	"context"
	"crypto/tls"
	confucius "github.com/Sanchous98/project-confucius-backend"
	"github.com/Sanchous98/project-confucius-backend/utils"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"time"
)

type server struct {
	config *Config
	server *http.Server
}

func NewServer(config utils.Config) confucius.Service {
	return &server{config.(*Config), &http.Server{Addr: ":80"}}
}

func (s *server) Init() error {
	return s.config.HydrateConfig()
}

func (s *server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = s.server.Shutdown(ctx)
	cancel()
}

func (s *server) Serve(router *mux.Router) error {
	// TODO: Test on real machine, because not working for localhost
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(os.Getenv("CERTS_PATH")),
		HostPolicy: autocert.HostWhitelist(s.config.Domain),
	}

	s.server.Handler = certManager.HTTPHandler(router)
	s.server.TLSConfig = &tls.Config{
		GetCertificate: certManager.GetCertificate,
	}

	go func() {
		log.Print("SERVER STARTED")
		log.Fatal(s.server.ListenAndServe())
	}()

	//log.Fatal(s.server.ListenAndServe())

	return nil
}
