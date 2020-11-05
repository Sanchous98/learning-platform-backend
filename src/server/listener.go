package server

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
	localGraphql "learning-platform/src/graphql"
	"log"
	"net/http"
	"os"
)

func Listen() {
	config := GraphQLConfig{}
	config.LoadConfig()
	router := mux.NewRouter()
	router.Handle("/api", handleApi())
	router.Handle("/graphiql", handleGraphiQL(os.Getenv("WORKING_PATH")+config.SchemaPath))
	serverConfig := ServerConfig{}
	serverConfig.LoadConfig()

	// TODO: Test on real machine, because not working for localhost
	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache(os.Getenv("CERTS_PATH")),
		// Put your domain here:
		HostPolicy: autocert.HostWhitelist(serverConfig.Domain),
	}

	server := &http.Server{
		Addr:    ":443",
		Handler: router,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go func() {
		log.Print("SERVER STARTED ON PORT 80")
		log.Fatal(http.ListenAndServe(":80", certManager.HTTPHandler(router)))
	}()

	log.Print("SERVER STARTED ON PORT 443")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func handleGraphiQL(schemaPath string) http.Handler {
	b, err := ioutil.ReadFile(schemaPath)

	if err != nil {
		panic("Schema file doesn't exist")
	}

	key := make([]byte, 32)
	_, _ = rand.Read(key)

	handlerFunc := handler.New(&handler.Config{
		Schema:     localGraphql.ResolveSchema(b),
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	return handlerFunc
}

func handleApi() http.Handler {
	router := mux.NewRouter()
	router.NewRoute().HandlerFunc(queryHandler)
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	log.Print(key)

	return csrf.Protect(key)(router)
}

func queryHandler(writer http.ResponseWriter, request *http.Request) {
	config := GraphQLConfig{}
	config.LoadConfig()
	b, err := ioutil.ReadFile(os.Getenv("WORKING_PATH") + config.SchemaPath)
	if err != nil {
		panic("Schema file doesn't exist")
	}

	query, err := ioutil.ReadAll(request.Body)

	if err != nil {
		log.Printf("Invalid query! Error: %v", err)
	}

	response := graphql.Do(graphql.Params{
		Schema:        *localGraphql.ResolveSchema(b),
		RequestString: string(query),
	})

	if len(response.Errors) > 0 {
		log.Printf("failed to execute graphql operation, errors: %+v", response.Errors)
	}

	writer.Header().Set("X-CSRF-Token", csrf.Token(request))
	jsonResponse, _ := json.Marshal(response)
	log.Print(writer.Write(jsonResponse))
}
