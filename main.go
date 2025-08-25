package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filePathRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filePathRoot)))
	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
