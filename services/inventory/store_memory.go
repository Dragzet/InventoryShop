package main

import "sync"

// ...existing code...

// InMemoryInventory — простая in-memory реализация, используется в тестах
type InMemoryInventory struct {
	mu     sync.Mutex
	items  map[int]*Item
	nextID int
}

func NewInventoryInMemory() *InMemoryInventory {
	return &InMemoryInventory{items: make(map[int]*Item), nextID: 1}
}

func (s *InMemoryInventory) List() []*Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]*Item, 0, len(s.items))
	for _, v := range s.items {
		res = append(res, v)
	}
	return res
}

func (s *InMemoryInventory) Get(id int) (*Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return it, nil
}

func (s *InMemoryInventory) Create(name string, qty int, price float64) *Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID
	s.nextID++
	it := &Item{ID: id, Name: name, Quantity: qty, Price: price}
	s.items[id] = it
	return it
}

func (s *InMemoryInventory) UpdateQuantity(id, delta int) (*Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	if it.Quantity+delta < 0 {
		return nil, ErrInsufficientStock
	}
	it.Quantity += delta
	return it, nil
}

var (
	ErrNotFound          = &customError{"not found"}
	ErrInsufficientStock = &customError{"insufficient stock"}
)

type customError struct{ msg string }

func (e *customError) Error() string { return e.msg }
