package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := "8002"
	if p := os.Getenv("ORDERS_PORT"); p != "" {
		port = p
	}
	invURL := "http://localhost:8001"
	if u := os.Getenv("INVENTORY_URL"); u != "" {
		invURL = u
	}

	router := NewRouter(invURL)
	log.Printf("Orders service listening on :%s (inventory: %s)", port, invURL)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
