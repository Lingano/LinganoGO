package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"LinganoGO/config"
	"LinganoGO/handlers"
	"LinganoGO/middleware" // Import the middleware package

	"github.com/gorilla/mux"
)

func main() {
	// Connect to MongoDB
	if err := config.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
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

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to LinganoGO")
	})

	// API Router
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Auth routes (public)
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	authRouter.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	// Authenticated API routes
	// All routes registered with authenticatedAPIRouter will use JWTMiddleware
	authenticatedAPIRouter := apiRouter.PathPrefix("/app").Subrouter() // Example prefix for authenticated app routes
	authenticatedAPIRouter.Use(middleware.JWTMiddleware)

	// User profile routes
	authenticatedAPIRouter.HandleFunc("/profile", handlers.GetUserProfile).Methods("GET")
	authenticatedAPIRouter.HandleFunc("/profile", handlers.UpdateUserProfile).Methods("PUT")

	// Saved words routes
	authenticatedAPIRouter.HandleFunc("/saved-words", handlers.AddSavedWord).Methods("POST")
	authenticatedAPIRouter.HandleFunc("/saved-words", handlers.GetSavedWords).Methods("GET")
	authenticatedAPIRouter.HandleFunc("/saved-words/{savedWordID}", handlers.DeleteSavedWord).Methods("DELETE")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}
}
