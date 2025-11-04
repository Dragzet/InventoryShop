package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// reserved represents a reserved quantity for rollback
type reserved struct{ id, qty int }

func NewRouter(invURL string) http.Handler {
	store := NewOrderStore()
	client := &http.Client{Timeout: 5 * time.Second}
	mux := http.NewServeMux()

	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, store.List())
		case http.MethodPost:
			var req struct {
				Items []struct {
					ID       int `json:"id"`
					Quantity int `json:"quantity"`
				} `json:"items"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
				return
			}

			var reservedList []reserved
			var orderItems []OrderItem
			var total float64

			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			for _, it := range req.Items {
				// fetch item details to get price
				getURL := fmt.Sprintf("%s/items/%d", strings.TrimRight(invURL, "/"), it.ID)
				reqGet, _ := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
				resp, err := client.Do(reqGet)
				if err != nil {
					rollbackInventory(ctx, client, invURL, reservedList)
					writeJSON(w, http.StatusBadGateway, map[string]string{"error": "failed to reach inventory"})
					return
				}
				if resp.StatusCode != http.StatusOK {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
					rollbackInventory(ctx, client, invURL, reservedList)
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": "item not found in inventory"})
					return
				}
				var invItem struct {
					ID       int     `json:"id"`
					Name     string  `json:"name"`
					Quantity int     `json:"quantity"`
					Price    float64 `json:"price"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&invItem); err != nil {
					resp.Body.Close()
					rollbackInventory(ctx, client, invURL, reservedList)
					writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "invalid inventory response"})
					return
				}
				resp.Body.Close()

				// send adjust request
				adjustURL := fmt.Sprintf("%s/items/%d/adjust", strings.TrimRight(invURL, "/"), it.ID)
				body, _ := json.Marshal(map[string]int{"delta": -it.Quantity})
				reqAdj, _ := http.NewRequestWithContext(ctx, http.MethodPost, adjustURL, bytes.NewReader(body))
				reqAdj.Header.Set("Content-Type", "application/json")
				resp2, err := client.Do(reqAdj)
				if err != nil {
					rollbackInventory(ctx, client, invURL, reservedList)
					writeJSON(w, http.StatusBadGateway, map[string]string{"error": "failed to adjust inventory"})
					return
				}
				io.Copy(io.Discard, resp2.Body)
				resp2.Body.Close()
				if resp2.StatusCode != http.StatusOK {
					rollbackInventory(ctx, client, invURL, reservedList)
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": "insufficient stock or invalid adjust"})
					return
				}

				// reserved ok
				reservedList = append(reservedList, reserved{id: it.ID, qty: it.Quantity})
				orderItems = append(orderItems, OrderItem{ItemID: it.ID, Name: invItem.Name, Quantity: it.Quantity, Price: invItem.Price})
				total += float64(it.Quantity) * invItem.Price
			}

			ord := store.Create(orderItems, total)
			writeJSON(w, http.StatusCreated, ord)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/orders/")
		if path == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		idStr := strings.Split(path, "/")[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		switch r.Method {
		case http.MethodGet:
			ord, err := store.Get(id)
			if err != nil {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
				return
			}
			writeJSON(w, http.StatusOK, ord)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// enable CORS and logging
	return loggingMiddleware(corsMiddleware(mux))
}

func rollbackInventory(ctx context.Context, client *http.Client, invURL string, reservedList []reserved) {
	for i := len(reservedList) - 1; i >= 0; i-- {
		r := reservedList[i]
		adjustURL := fmt.Sprintf("%s/items/%d/adjust", strings.TrimRight(invURL, "/"), r.id)
		body, _ := json.Marshal(map[string]int{"delta": r.qty})
		reqAdj, _ := http.NewRequestWithContext(ctx, http.MethodPost, adjustURL, bytes.NewReader(body))
		reqAdj.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(reqAdj)
		if err != nil {
			log.Printf("rollback failed for item %d: %v", r.id, err)
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware adds permissive CORS headers for local development
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
