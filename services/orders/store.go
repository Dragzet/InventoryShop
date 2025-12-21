package main

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"time"
)

// OrderItem represents item in an order
type OrderItem struct {
	ItemID   int     `json:"item_id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// Order represents a customer's order
type Order struct {
	ID      int         `json:"id"`
	Items   []OrderItem `json:"items"`
	Total   float64     `json:"total"`
	Created int64       `json:"created_unix"`
}

// OrderStore is a Postgres-backed store for orders
type OrderStore struct {
	db *sql.DB
}

func NewOrderStore(db *sql.DB) *OrderStore {
	// create tables
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		total NUMERIC NOT NULL,
		created_unix BIGINT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS order_items (
		id SERIAL PRIMARY KEY,
		order_id INT REFERENCES orders(id) ON DELETE CASCADE,
		item_id INT NOT NULL,
		name TEXT NOT NULL,
		quantity INT NOT NULL,
		price NUMERIC NOT NULL
	);
	`)
	if err != nil {
		panic(err)
	}
	return &OrderStore{db: db}
}

func (s *OrderStore) Create(items []OrderItem, total float64) *Order {
	var orderID int
	// transactional insert
	tx, err := s.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()
	if err := tx.QueryRow("INSERT INTO orders (total, created_unix) VALUES ($1, $2) RETURNING id", total, nowUnix()).Scan(&orderID); err != nil {
		panic(err)
	}
	for _, it := range items {
		_, err := tx.Exec("INSERT INTO order_items (order_id, item_id, name, quantity, price) VALUES ($1,$2,$3,$4,$5)", orderID, it.ItemID, it.Name, it.Quantity, it.Price)
		if err != nil {
			panic(err)
		}
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}
	return &Order{ID: orderID, Items: items, Total: total, Created: nowUnix()}
}

func (s *OrderStore) Get(id int) (*Order, error) {
	var o Order
	row := s.db.QueryRow("SELECT id, total, created_unix FROM orders WHERE id=$1", id)
	if err := row.Scan(&o.ID, &o.Total, &o.Created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	rows, err := s.db.Query("SELECT item_id, name, quantity, price FROM order_items WHERE order_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]OrderItem, 0)
	for rows.Next() {
		var it OrderItem
		if err := rows.Scan(&it.ItemID, &it.Name, &it.Quantity, &it.Price); err != nil {
			continue
		}
		items = append(items, it)
	}
	o.Items = items
	return &o, nil
}

func (s *OrderStore) List() []*Order {
	rows, err := s.db.Query("SELECT id, total, created_unix FROM orders ORDER BY id DESC")
	if err != nil {
		return []*Order{}
	}
	defer rows.Close()
	res := make([]*Order, 0)
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.Total, &o.Created); err != nil {
			continue
		}
		// load items
		itRows, err := s.db.Query("SELECT item_id, name, quantity, price FROM order_items WHERE order_id=$1", o.ID)
		if err == nil {
			for itRows.Next() {
				var it OrderItem
				if err := itRows.Scan(&it.ItemID, &it.Name, &it.Quantity, &it.Price); err == nil {
					o.Items = append(o.Items, it)
				}
			}
			itRows.Close()
		}
		res = append(res, &o)
	}
	return res
}

func nowUnix() int64 { return time.Now().Unix() }
