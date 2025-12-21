package main

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

// Item represents a product in inventory
type Item struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// Inventory is a Postgres-backed store
type Inventory struct {
	db *sql.DB
}

// NewInventory initializes store and ensures table exists
func NewInventory(db *sql.DB) *Inventory {
	// create table if not exists
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		quantity INT NOT NULL,
		price NUMERIC NOT NULL
	)
	`)
	if err != nil {
		panic(err)
	}
	return &Inventory{db: db}
}

func (s *Inventory) List() []*Item {
	rows, err := s.db.Query("SELECT id, name, quantity, price FROM items")
	if err != nil {
		return []*Item{}
	}
	defer rows.Close()
	res := make([]*Item, 0)
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.Name, &it.Quantity, &it.Price); err != nil {
			continue
		}
		res = append(res, &it)
	}
	return res
}

func (s *Inventory) Get(id int) (*Item, error) {
	var it Item
	row := s.db.QueryRow("SELECT id, name, quantity, price FROM items WHERE id=$1", id)
	if err := row.Scan(&it.ID, &it.Name, &it.Quantity, &it.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &it, nil
}

func (s *Inventory) Create(name string, qty int, price float64) *Item {
	var id int
	err := s.db.QueryRow("INSERT INTO items (name, quantity, price) VALUES ($1,$2,$3) RETURNING id", name, qty, price).Scan(&id)
	if err != nil {
		panic(err)
	}
	return &Item{ID: id, Name: name, Quantity: qty, Price: price}
}

func (s *Inventory) UpdateQuantity(id, delta int) (*Item, error) {
	// Try to update only when resulting quantity >= 0 and return the row
	var it Item
	row := s.db.QueryRow(`
	UPDATE items SET quantity = quantity + $1
	WHERE id = $2 AND (quantity + $1) >= 0
	RETURNING id, name, quantity, price
	`, delta, id)
	if err := row.Scan(&it.ID, &it.Name, &it.Quantity, &it.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found or insufficient stock")
		}
		return nil, err
	}
	return &it, nil
}
