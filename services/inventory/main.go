package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := "8001"
	if p := os.Getenv("INVENTORY_PORT"); p != "" {
		port = p
	}

	router := NewRouter()
	log.Printf("Inventory service listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
