package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"LinganoGO/config"
	"LinganoGO/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const defaultPort = "8080"

func main() {
	// Connect to PostgreSQL
	if err := config.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	// Defer disconnection
	defer config.DisconnectDB()

	// Set up signal handling for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopChan
		log.Println("Shutting down server...")
		config.DisconnectDB() // Ensure DB is disconnected on shutdown
		os.Exit(0)
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(middleware.Logger)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

