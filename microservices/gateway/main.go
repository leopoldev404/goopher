package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	var mux = http.NewServeMux()
	var handler = NewHandler()
	handler.RegisterRoutes(mux)

	log.Print("Started Gateway Service! 🚀")
	if err := http.ListenAndServe(os.Getenv("GATEWAY_SERVICE_PORT"), mux); err != nil {
		log.Fatal("Failed Starting Http Gateway Server! 😱", err)
	}
}
