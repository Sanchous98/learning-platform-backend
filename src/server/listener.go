package server

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
)

var middlewares = []func(handler http.Handler) *http.Handler{
	LoggerMiddleware,
	AuthenticationMiddleware,
	AuthorizationMiddleware,
	CSRFMiddleware,
}

func Listen() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", queryHandler)
	serverConfig := ServerConfig{}
	serverConfig.LoadConfig()

	// TODO: Test on real machine, because not working for localhost
	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("cert-cache"),
		// Put your domain here:
		HostPolicy: autocert.HostWhitelist(serverConfig.Domain),
	}

	server := &http.Server{
		Addr:    ":443",
		Handler: getGlobalMiddlewares(mux),
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go func() {
		log.Print("SERVER STARTED ON PORT 80")
		log.Fatal(http.ListenAndServe(":80", certManager.HTTPHandler(nil)))
	}()

	log.Print("SERVER STARTED ON PORT 443")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world"))
	// TODO: Handle GraphQL query
}

func getGlobalMiddlewares(handle http.Handler) http.Handler {
	var handled http.Handler
	for _, middleware := range middlewares {
		handled = *middleware(handle)
	}

	return handled
}
