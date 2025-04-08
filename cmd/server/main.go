package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/yourusername/go-saas-template/internal/auth"
)

// loggingMiddleware logs information about incoming HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("incoming request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		next.ServeHTTP(w, r)

		log.Printf("✅ completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize PocketBase data directory with better environment handling
	pbDataDir := os.Getenv("PB_DATA_DIR")
	if pbDataDir == "" {
		// Default to ./pb_data for development
		pbDataDir = "./pb_data"

		// Create the data directory if it doesn't exist (not strictly necessary, but explicit)
		if err := os.MkdirAll(pbDataDir, os.ModePerm); err != nil {
			log.Printf("⚠️ Warning: Could not create data directory: %v", err)
		}
	}

	log.Printf("📁 Using PocketBase data directory: %s", filepath.Clean(pbDataDir))

	// Check if this is a fresh installation
	if _, err := os.Stat(filepath.Join(pbDataDir, "data.db")); os.IsNotExist(err) {
		log.Println("🆕 Fresh installation detected, initializing default data...")
		// Initialize your default collections, users, etc.
	}

	// Create a new PocketBase app
	pb := pocketbase.New()

	// Set the app as the global PbClient
	auth.PbClient = pb

	// Add a hook to log when PocketBase is initialized
	// Use OnServe().BindFunc() instead of OnServe().Add() as per documentation
	pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		log.Println("🔐 PocketBase initialized successfully")
		return e.Next()
	})

	// Initialize PocketBase (without starting the server)
	if err := pb.Bootstrap(); err != nil {
		log.Fatalf("❌ Failed to initialize PocketBase: %v", err)
	}

	r := mux.NewRouter()

	// Apply logging middleware to all routes
	r.Use(loggingMiddleware)

	// Auth routes - these don't require authentication
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", auth.LoginHandler).Methods("GET", "POST")
	authRouter.HandleFunc("/register", auth.RegisterHandler).Methods("GET", "POST")
	authRouter.HandleFunc("/logout", auth.LogoutHandler).Methods("GET")
	authRouter.HandleFunc("/forgot-password", auth.ForgotPasswordHandler).Methods("GET", "POST")
	authRouter.HandleFunc("/reset-password", auth.ResetPasswordHandler).Methods("GET", "POST")

	// PocketBase auth API forwarding
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/auth/{action}", auth.PocketBaseAuthHandler).Methods("POST")

	// Protected routes - require authentication
	protectedRouter := r.PathPrefix("/").Subrouter()
	protectedRouter.Use(auth.AuthMiddleware)

	// Dashboard/Home page (protected)
	protectedRouter.HandleFunc("/", auth.HomeRenderer)

	log.Printf("🚀 Starting HTTP server on port %s (http://localhost:%s)", port, port)
	log.Println("Press Ctrl+C to stop the server")

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("❌ Failed to start web server on port %s: %v", port, err)
	}
}
