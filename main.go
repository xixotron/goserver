package main

import (
	"log"
	"net/http"
)

const serverPort = "8080"

func main() {

	serverMux := http.NewServeMux()

	httpServer := http.Server{
		Addr:    ":" + serverPort,
		Handler: serverMux,
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
