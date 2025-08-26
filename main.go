package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const port = "8080"
	const filePathRoot = "."

	appCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	fsHandler := appCfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))
	mux.Handle("/app/", http.StripPrefix("/app/", fsHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handleChirpsValidate)

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
	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
