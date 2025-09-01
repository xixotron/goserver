package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/xixotron/goserver/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}
	dbQueries := database.New(db)

	platform := os.Getenv("PLATFORM")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY must be set")
	}

	appCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecretKey,
	}

	mux := http.NewServeMux()

	fsHandler := appCfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))
	mux.Handle("/app/", http.StripPrefix("/app/", fsHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/users", appCfg.handleCreateUser)
	mux.HandleFunc("POST /api/login", appCfg.handleLogin)
	mux.HandleFunc("POST /api/revoke", appCfg.handleRevokeRefreshToken)
	mux.HandleFunc("POST /api/refresh", appCfg.handleRefreshToken)

	mux.HandleFunc("POST /api/chirps", appCfg.handlePostChirp)
	mux.HandleFunc("GET /api/chirps", appCfg.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", appCfg.handleGetChirp)

	mux.Handle("GET /admin/metrics", appCfg.handlerMetrics())
	mux.HandleFunc("POST /admin/reset", appCfg.handlerMetricsReset)
	// Both Handle and Handle Func take a pattern string and then regisger a handler to serve queryes on that pattern
	// Handle takes a http.handler wich implements ServeHTTP(w http.ResponseWriter, r *http.Request) as the handler
	// HandleFunc takes a funcion with signature func(w http.ResponseWriter, r *http.Request) and Registers it as a Handler

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}

	log.Println("Bye!")
}
