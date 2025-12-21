package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	port := "8002"
	if p := os.Getenv("ORDERS_PORT"); p != "" {
		port = p
	}
	invURL := "http://inventory:8001"
	if u := os.Getenv("INVENTORY_URL"); u != "" {
		invURL = u
	}

	pg := os.Getenv("ORDERS_DATABASE_URL")
	if pg == "" {
		pg = "postgres://postgres:postgres@orders-db:5432/orders?sslmode=disable"
	}
	db, err := sql.Open("postgres", pg)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping orders db: %v", err)
	}

	store := NewOrderStore(db)
	router := NewRouter(store, invURL)
	log.Printf("Orders service listening on :%s (inventory: %s)", port, invURL)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
