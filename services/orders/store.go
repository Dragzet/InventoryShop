package main

import (
	"errors"
	"sync"
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

// OrderStore is an in-memory store for orders
type OrderStore struct {
	mu     sync.Mutex
	orders map[int]*Order
	nextID int
}

func NewOrderStore() *OrderStore {
	return &OrderStore{orders: make(map[int]*Order), nextID: 1}
}

func (s *OrderStore) Create(items []OrderItem, total float64) *Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID
	s.nextID++
	ord := &Order{ID: id, Items: items, Total: total, Created: nowUnix()}
	s.orders[id] = ord
	return ord
}

func (s *OrderStore) Get(id int) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return o, nil
}

func (s *OrderStore) List() []*Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]*Order, 0, len(s.orders))
	for _, v := range s.orders {
		res = append(res, v)
	}
	return res
}

func nowUnix() int64 { return time.Now().Unix() }
