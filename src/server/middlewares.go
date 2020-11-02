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

func AuthenticationMiddleware(next http.Handler) *Authentication {
	return &Authentication{Middleware{next}}
}

func AuthorizationMiddleware(next http.Handler) *Authorization {
	return &Authorization{Middleware{next}}
}

func CSRFMiddleware(next http.Handler) *CSRF {
	return &CSRF{Middleware{next}}
}

func LoggerMiddleware(next http.Handler) *Logger {
	return &Logger{Middleware{next}}
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
