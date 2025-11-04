package main

import (
	"errors"
	"sync"
)

// Item represents a product in inventory
type Item struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// Inventory is a simple in-memory store
type Inventory struct {
	mu     sync.Mutex
	items  map[int]*Item
	nextID int
}

func NewInventory() *Inventory {
	return &Inventory{items: make(map[int]*Item), nextID: 1}
}

func (s *Inventory) List() []*Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]*Item, 0, len(s.items))
	for _, v := range s.items {
		res = append(res, v)
	}
	return res
}

func (s *Inventory) Get(id int) (*Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return it, nil
}

func (s *Inventory) Create(name string, qty int, price float64) *Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID
	s.nextID++
	it := &Item{ID: id, Name: name, Quantity: qty, Price: price}
	s.items[id] = it
	return it
}

func (s *Inventory) UpdateQuantity(id, delta int) (*Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	if it.Quantity+delta < 0 {
		return nil, errors.New("insufficient stock")
	}
	it.Quantity += delta
	return it, nil
}
