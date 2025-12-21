package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func NewRouter(store *Inventory) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			list := store.List()
			writeJSON(w, http.StatusOK, list)
		case http.MethodPost:
			var req struct {
				Name     string  `json:"name"`
				Quantity int     `json:"quantity"`
				Price    float64 `json:"price"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
				return
			}
			it := store.Create(req.Name, req.Quantity, req.Price)
			writeJSON(w, http.StatusCreated, it)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		// expected: /items/{id} or /items/{id}/adjust
		path := strings.TrimPrefix(r.URL.Path, "/items/")
		if path == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		parts := strings.Split(path, "/")
		idStr := parts[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}

		if len(parts) == 1 {
			switch r.Method {
			case http.MethodGet:
				it, err := store.Get(id)
				if err != nil {
					writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
					return
				}
				writeJSON(w, http.StatusOK, it)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}

		// path like {id}/adjust
		if parts[1] == "adjust" && r.Method == http.MethodPost {
			// read delta from JSON body {"delta": -2}
			var req struct {
				Delta int `json:"delta"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
				return
			}
			it, err := store.UpdateQuantity(id, req.Delta)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, it)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	return loggingMiddleware(corsMiddleware(mux))
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
		origin := r.Header.Get("Origin")
		if origin == "" {
			// no origin, proceed normally
			next.ServeHTTP(w, r)
			return
		}

		// echo back the Origin to allow cookies/auth if needed
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			// preflight
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
