package main

import (
	confucius "github.com/Sanchous98/project-confucius-backend"
	"github.com/Sanchous98/project-confucius-backend/service/graphql"
	"github.com/Sanchous98/project-confucius-backend/service/server"
	"github.com/gorilla/mux"
	"log"
	"os"
)

var Container = confucius.NewContainer()

func main() {
	bootstrap()
	router := mux.NewRouter()

	if err := Container.Init(); err != nil {
		log.Fatal(err)
	}

	if err := Container.Serve(router); err != nil {
		log.Fatal(err)
	}
}

func bootstrap() {
	// Set working dir
	workDir, err := os.Getwd()
	err = os.Setenv("WORKING_PATH", workDir)

	if err != nil {
		panic("No working dir")
	}

	// Set config path
	err = os.Setenv("CONFIG_PATH", workDir+"/config")

	if err != nil {
		panic("Cannot set config path")
	}

	if _, err := os.Stat(os.Getenv("CONFIG_PATH")); os.IsNotExist(err) {
		err = os.Mkdir(os.Getenv("CONFIG_PATH"), 0755)

		if err != nil {
			panic("Config path doesn't exists and cannot be created")
		}
	}

	// Set certificates path
	err = os.Setenv("CERTS_PATH", workDir+"/certs")

	if err != nil {
		panic("Cannot set certs path")
	}

	if _, err := os.Stat(os.Getenv("CERTS_PATH")); os.IsNotExist(err) {
		err = os.Mkdir(os.Getenv("CERTS_PATH"), 0755)

		if err != nil {
			panic("Certs path doesn't exists and cannot be created")
		}
	}

	Container.SetMainService("server", server.NewServer(&server.Config{}))
	Container.Set("graphql", graphql.NewService(&graphql.Config{}))
}
