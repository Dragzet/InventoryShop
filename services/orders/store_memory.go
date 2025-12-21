package main

import (
	"errors"
	"sync"
	"time"
)

// In-memory store for orders (used in tests)
type OrderStoreInMemory struct {
	mu     sync.Mutex
	orders map[int]*Order
	nextID int
}

func NewOrderStoreInMemory() *OrderStoreInMemory {
	return &OrderStoreInMemory{orders: make(map[int]*Order), nextID: 1}
}

func (s *OrderStoreInMemory) Create(items []OrderItem, total float64) *Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID
	s.nextID++
	ord := &Order{ID: id, Items: items, Total: total, Created: time.Now().Unix()}
	s.orders[id] = ord
	return ord
}

func (s *OrderStoreInMemory) Get(id int) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return o, nil
}

func (s *OrderStoreInMemory) List() []*Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]*Order, 0, len(s.orders))
	for _, o := range s.orders {
		res = append(res, o)
	}
	return res
}
