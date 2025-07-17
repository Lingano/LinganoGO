package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"LinganoGO/config"
	"LinganoGO/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fatih/color"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
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
    ░  ░ ░           ░       ░       ░  ░         ░     ░ ░  `)
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
	
	// Add CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver()}))
	
	// Enable introspection for GraphQL Playground docs
	srv.Use(extension.Introspection{})
	
	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/query", srv)
	router.Handle("/graphql", srv) // Main GraphQL endpoint
	
	// Add a health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "message": "LinganoGO GraphQL Server is running"}`))
	})
	
	color.Red(horizontalLine)
	color.New(color.FgYellow).Print("The port is set to: ")
	color.New(color.FgGreen).Print(port)
	color.New(color.FgYellow).Print("\n")
	color.Yellow("You can change it in the .env file")
	color.Red(horizontalLine)
	color.Green("🚀 Server ready at http://localhost:%s/", port)
	color.Green("📊 GraphQL Playground: http://localhost:%s/", port)
	color.Green("🔍 Health Check: http://localhost:%s/health", port)
	color.Red(horizontalLine)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
