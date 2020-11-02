package server

import (
	"net/http"
)

type Middleware struct {
	handler http.Handler
}

type Authentication struct {
	Middleware
}

type Authorization struct {
	Middleware
}

type CSRF struct {
	excludedPaths []string
	useCookie     bool
	Middleware
}

type Logger struct {
	Middleware
}

func AuthenticationMiddleware(next http.Handler) *http.Handler {
	var handler http.Handler = Authentication{Middleware{next}}
	return &handler
}

func AuthorizationMiddleware(next http.Handler) *http.Handler {
	var handler http.Handler = Authorization{Middleware{next}}
	return &handler
}

func CSRFMiddleware(next http.Handler) *http.Handler {
	serverConfig := ServerConfig{}
	serverConfig.LoadConfig()
	var handler http.Handler = CSRF{serverConfig.CSRF.ExcludedPaths, serverConfig.CSRF.UseCookies, Middleware{next}}
	return &handler
}

func LoggerMiddleware(next http.Handler) *http.Handler {
	var handler http.Handler = Logger{Middleware{next}}
	return &handler
}

func (auth Authentication) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auth.handler.ServeHTTP(w, r)
}

func (auth Authorization) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auth.handler.ServeHTTP(w, r)
}

func (csrf CSRF) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	csrf.handler.ServeHTTP(w, r)
}

func (log Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.handler.ServeHTTP(w, r)
}
