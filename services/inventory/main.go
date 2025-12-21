package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	port := "8001"
	if p := os.Getenv("INVENTORY_PORT"); p != "" {
		port = p
	}

	pg := os.Getenv("INVENTORY_DATABASE_URL")
	if pg == "" {
		pg = "postgres://postgres:postgres@inventory-db:5432/inventory?sslmode=disable"
	}
	db, err := sql.Open("postgres", pg)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// try ping
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	store := NewInventory(db)

	// seed demo data if empty
	rows := store.List()
	if len(rows) == 0 {
		store.Create("Толстовка", 100, 19.99)
		store.Create("Футболка", 50, 7.5)
	}

	router := NewRouter(store)
	log.Printf("Inventory service listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
