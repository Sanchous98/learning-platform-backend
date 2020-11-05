package server

import (
	"log"
	"net/http"
)

func Listen() {
	http.HandleFunc("/", queryHandler)
	log.Print("SERVER STARTED")
	// TODO: Secure connection. Use "Let's encrypt" to get TLS certificate
	log.Fatal(http.ListenAndServe(":80", nil))
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world"))
	// TODO: Handle GraphQL query
}
