package main

import (
	"log"
	"net/http"

	"github.com/mdcurran/exchange-rate/internal/pkg/api"
)

func main() {
	s := api.NewServer()

	log.Printf("Starting HTTP server\n")
	err := http.ListenAndServe(":8080", s.Router)
	if err != nil {
		log.Fatalf("Unable to initialise API server - %v", err)
	}
}
