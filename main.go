package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/wader/gormstore/v2"
)

var (
	store sessions.Store
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	initDB()

	// Initialize session store with gormstore
	store = gormstore.New(db, []byte(os.Getenv("SESSION_SECRET")))
	store.(*gormstore.Store).SessionOpts = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   true, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
	}

	// Cleanup old sessions
	go store.(*gormstore.Store).PeriodicCleanup(1*time.Hour, nil)

	// Initialize router
	r := mux.NewRouter()

	// Middleware
	r.Use(loggingMiddleware)

	// Public routes
	r.HandleFunc("/", handleHome).Methods("GET")
	authRouter(r.PathPrefix("/").Subrouter())

	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(authMiddleware)

	// Get port from environment
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Welcome to PolyglotWeb API",
		"status":  "running",
	})
}
