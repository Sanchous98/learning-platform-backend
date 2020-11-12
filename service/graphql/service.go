package graphql

import (
	"encoding/json"
	"fmt"
	confucius "github.com/Sanchous98/project-confucius-backend"
	"github.com/Sanchous98/project-confucius-backend/utils"
	tools "github.com/bhoriuchi/graphql-go-tools"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
)

// graphQL service is used to add some http entries that serve GraphQL queries
type graphQL struct {
	config *Config
}

func NewService(config utils.Config) confucius.Service {
	return &graphQL{config.(*Config)}
}

func (g *graphQL) Serve(router *mux.Router) error {
	router.Handle("/api", csrfMiddleware(*g))
	router.Handle("/graphiql", handleGraphiQL(os.Getenv("WORKING_PATH")+g.config.SchemaPath))
	return nil
}

// Only for implementing an interface. Nothing to stop
func (g *graphQL) Stop() {}

func (g *graphQL) Init() error {
	return g.config.HydrateConfig()
}

// csrfMiddleware is a middleware that prevents server from XSRF and XSS attacks
func csrfMiddleware(g graphQL) http.Handler {
	router := mux.NewRouter()
	router.NewRoute().HandlerFunc(g.queryHandler)
	key := make([]byte, 32)
	_, _ = rand.Read(key)

	return csrf.Protect(key)(router)
}

// resolveSchema initializes GraphQL server configuration, based on schema file and predefined resolvers and directives
func resolveSchema(schemaContent []byte) *graphql.Schema {
	schema, err := tools.MakeExecutableSchema(tools.ExecutableSchema{
		TypeDefs: string(schemaContent),
		SchemaDirectives: tools.SchemaDirectiveVisitorMap{
			"IsGranted": &tools.SchemaDirectiveVisitor{
				VisitFieldDefinition: IsGranted,
			},
		},
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to parse schema, error: %v", err))
	}

	return &schema
}

// handleGraphiQL provides GraphiQL playground
func handleGraphiQL(schemaPath string) http.Handler {
	b, err := ioutil.ReadFile(schemaPath)

	if err != nil {
		panic("Schema file doesn't exist")
	}

	return handler.New(&handler.Config{
		Schema:     resolveSchema(b),
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})
}

// queryHandler is a function that handles GraphQL queries
func (g graphQL) queryHandler(writer http.ResponseWriter, request *http.Request) {
	b, err := ioutil.ReadFile(os.Getenv("WORKING_PATH") + g.config.SchemaPath)
	if err != nil {
		panic("Schema file doesn't exist")
	}

	query, err := ioutil.ReadAll(request.Body)

	if err != nil {
		log.Printf("Invalid query! Error: %v", err)
	}

	response := graphql.Do(graphql.Params{
		Schema:        *resolveSchema(b),
		RequestString: string(query),
	})

	if len(response.Errors) > 0 {
		log.Printf("Failed to execute graphql operation, errors: %+v", response.Errors)
	}

	writer.Header().Set("X-CSRF-Token", csrf.Token(request))
	jsonResponse, _ := json.Marshal(response)
	log.Print(writer.Write(jsonResponse))
}
