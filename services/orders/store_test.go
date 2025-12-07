package main

import (
	"testing"
	"time"
)

func almostEqualFloat(a, b float64) bool {
	if a == b {
		return true
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < 1e-9
}

func TestOrderStore_CreateGetList(t *testing.T) {
	s := NewOrderStore()
	items := []OrderItem{{ItemID: 1, Name: "apple", Quantity: 2, Price: 1.5}}
	ord := s.Create(items, 3.0)
	if ord.ID != 1 {
		t.Fatalf("expected id 1, got %d", ord.ID)
	}
	if !almostEqualFloat(ord.Total, 3.0) {
		t.Fatalf("expected total 3.0, got %v", ord.Total)
	}
	if ord.Created <= 0 {
		t.Fatalf("expected positive created timestamp, got %d", ord.Created)
	}

	// small time drift check
	if time.Unix(ord.Created, 0).After(time.Now().Add(2 * time.Second)) {
		t.Fatalf("created timestamp is in the future: %v", ord.Created)
	}

	// Get
	got, err := s.Get(ord.ID)
	if err != nil {
		t.Fatalf("unexpected error from Get: %v", err)
	}
	if got.ID != ord.ID || len(got.Items) != 1 {
		t.Fatalf("unexpected order data: %+v", got)
	}

	// List
	list := s.List()
	if len(list) != 1 {
		t.Fatalf("expected list len 1, got %d", len(list))
	}

	// Get non-existing
	_, err = s.Get(999)
	if err == nil {
		t.Fatalf("expected error getting non-existing order, got nil")
	}
}
