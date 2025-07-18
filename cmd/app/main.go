package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	httphandler "whatsapp-parser/internal/delivery/http"
	"whatsapp-parser/internal/repository"
	"whatsapp-parser/internal/usecase"
)

func main() {
	// Create storage directory
	storageDir := filepath.Join(".", "storage", "sessions")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	// Initialize repository
	sessionRepo, err := repository.NewSessionRepository(storageDir)
	if err != nil {
		log.Fatalf("Failed to create session repository: %v", err)
	}

	// Initialize use case
	sessionUseCase, err := usecase.NewSessionUseCase(sessionRepo)
	if err != nil {
		log.Fatalf("Failed to create session use case: %v", err)
	}

	// Initialize HTTP handler
	h := httphandler.NewHandler(sessionUseCase)

	// Create router
	r := mux.NewRouter()
	h.RegisterRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
} 