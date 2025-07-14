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
const horizontalLine = "================================================================"

func displayTitle() {
	color.New(color.FgCyan, color.Bold).Println(`
 ██▓     ██▓ ███▄    █   ▄████  ▄▄▄       ███▄    █  ▒█████  
▓██▒    ▓██▒ ██ ▀█   █  ██▒ ▀█▒▒████▄     ██ ▀█   █ ▒██▒  ██▒
▒██░    ▒██▒▓██  ▀█ ██▒▒██░▄▄▄░▒██  ▀█▄  ▓██  ▀█ ██▒▒██░  ██▒
▒██░    ░██░▓██▒  ▐▌██▒░▓█  ██▓░██▄▄▄▄██ ▓██▒  ▐▌██▒▒██   ██░
░██████▒░██░▒██░   ▓██░░▒▓███▀▒ ▓█   ▓██▒▒██░   ▓██░░ ████▓▒░
░ ▒░▓  ░░▓  ░ ▒░   ▒ ▒  ░▒   ▒  ▒▒   ▓▒█░░ ▒░   ▒ ▒ ░ ▒░▒░▒░ 
░ ░ ▒  ░ ▒ ░░ ░░   ░ ▒░  ░   ░   ▒   ▒▒ ░░ ░░   ░ ▒░  ░ ▒ ▒░ 
  ░ ░    ▒ ░   ░   ░ ░ ░ ░   ░   ░   ▒      ░   ░ ░ ░ ░ ░ ▒  
    ░  ░ ░           ░       ░       ░  ░         ░     ░ ░  
`)
	color.New(color.FgMagenta, color.Bold).Println("                    GraphQL API Server")
	color.New(color.FgYellow).Println("                      Version 1.0.0")
}

func main() {
	// Display title
	displayTitle()

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
	color.Red(horizontalLine)
	color.New(color.FgYellow).Print("The port is set to: ")
	color.New(color.FgGreen).Print(port)
	color.New(color.FgYellow).Print("\n")
	color.Yellow("You can change it in the .env file")
	color.Red(horizontalLine)
	color.Green("🚀 Server ready at http://localhost:%s/", port)
	color.Red(horizontalLine)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
