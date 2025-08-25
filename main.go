package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filePathRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot))))
	mux.Handle("/healthz", handlerReadiness())

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
