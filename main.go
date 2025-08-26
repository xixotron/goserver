package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	const port = "8080"
	const filePathRoot = "."

	appCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", appCfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))))
	mux.Handle("/healthz", handlerReadiness())
	mux.Handle("/metrics", appCfg.handlerStatistics())
	mux.Handle("/reset", appCfg.handlerReset())

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

func handlerReadiness() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerStatistics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hits: %v", cfg.fileServerHits.Load())
	})
}

func (cfg *apiConfig) handlerReset() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Store(0)

		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hits: %v", cfg.fileServerHits.Load())
	})
}
