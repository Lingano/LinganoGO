package main

import (
	"log"
	"net/http"
	"os"

	"LinganoGO/config"
	"LinganoGO/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fatih/color"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const defaultPort = "8081"

func main() {
	// Connect to PostgreSQL (keep for existing migrations)
	if err := config.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer config.DisconnectDB()

	// Connect to PostgreSQL with Ent
	if err := config.ConnectEntDB(); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL with Ent: %v", err)
	}
	defer config.DisconnectEntDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	color.Green("🚀 Server ready at http://localhost:%s/", port)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
