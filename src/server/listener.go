package server

import (
	"crypto/tls"
	"encoding/json"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/api", queryHandler)
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
		Handler: getGlobalMiddlewares(mux),
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	log.Print("SERVER STARTED ON PORT 80")
	go log.Fatal(http.ListenAndServe(":80", certManager.HTTPHandler(nil)))

	log.Print("SERVER STARTED ON PORT 443")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func GraphiQL() {
	config := GraphQLConfig{}
	config.LoadConfig()
	b, err := ioutil.ReadFile(os.Getenv("WORKING_PATH") + config.SchemaPath)

	if err != nil {
		panic("Schema file doesn't exist")
	}

	http.Handle("/graphiql", handler.New(&handler.Config{
		Schema:     localGraphql.ResolveSchema(b),
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	}))
	_ = http.ListenAndServe(":8080", nil)
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
		VariableValues: map[string]interface{}{
			"course": map[string]interface{}{
				"id": "id",
			},
		},
	})

	if len(response.Errors) > 0 {
		log.Printf("failed to execute graphql operation, errors: %+v", response.Errors)
	}

	jsonResponse, _ := json.Marshal(response)
	log.Print(writer.Write(jsonResponse))
}

func getGlobalMiddlewares(handle http.Handler) http.Handler {
	var handled http.Handler
	middlewares := []func(http.Handler) *http.Handler{
		LoggerMiddleware,
		AuthenticationMiddleware,
		AuthorizationMiddleware,
		CSRFMiddleware,
	}

	for _, middleware := range middlewares {
		handled = *middleware(handle)
	}

	return handled
}
