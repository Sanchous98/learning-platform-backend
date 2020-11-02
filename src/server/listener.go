package server

import (
	"log"
	"net/http"
)

var middlewares = []func(handle http.Handler) http.Handler{
	LoggerMiddleware,
}

func Listen() {
	server := http.NewServeMux()
	server.HandleFunc("/", queryHandler)
	wrappedServer := middleware(server)
	log.Print("SERVER STARTED")
	// TODO: Secure connection. Use "Let's encrypt" to get TLS certificate
	log.Fatal(http.ListenAndServe(":80", wrappedServer))
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world"))
	// TODO: Handle GraphQL query
}

func middleware(handle http.Handler) http.Handler {
	var handled http.Handler = LoggerMiddleware(handle)
	handled = AuthenticationMiddleware(handle)
	handled = AuthorizationMiddleware(handled)
	handled = CSRFMiddleware(handled)

	return handled
}
