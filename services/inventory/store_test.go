package main

import (
	"testing"
)

func TestInventory_CreateGetList_UpdateQuantity(t *testing.T) {
	s := NewInventoryInMemory()
	it1 := s.Create("apple", 10, 1.5)
	if it1.ID != 1 {
		t.Fatalf("expected id 1, got %d", it1.ID)
	}
	it2 := s.Create("banana", 5, 2.0)
	if it2.ID != 2 {
		t.Fatalf("expected id 2, got %d", it2.ID)
	}

	// Get existing
	got, err := s.Get(it1.ID)
	if err != nil {
		t.Fatalf("unexpected error from Get: %v", err)
	}
	if got.Name != "apple" || got.Quantity != 10 {
		t.Fatalf("unexpected item data: %+v", got)
	}

	// List
	list := s.List()
	if len(list) != 2 {
		t.Fatalf("expected list length 2, got %d", len(list))
	}

	// UpdateQuantity success
	updated, err := s.UpdateQuantity(it1.ID, -3)
	if err != nil {
		t.Fatalf("unexpected error from UpdateQuantity: %v", err)
	}
	if updated.Quantity != 7 {
		t.Fatalf("expected quantity 7, got %d", updated.Quantity)
	}

	// UpdateQuantity insufficient stock
	_, err = s.UpdateQuantity(it2.ID, -10)
	if err == nil {
		t.Fatalf("expected error for insufficient stock, got nil")
	}

	// Get non-existing
	_, err = s.Get(9999)
	if err == nil {
		t.Fatalf("expected error when getting non-existing item, got nil")
	}
}
