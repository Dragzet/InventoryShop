package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch cmd {
	case "list-items":
		listItems()
	case "create-item":
		if len(os.Args) < 5 {
			fmt.Println("usage: create-item NAME QUANTITY PRICE")
			os.Exit(1)
		}
		qty, _ := strconv.Atoi(os.Args[3])
		price, _ := strconv.ParseFloat(os.Args[4], 64)
		createItem(os.Args[2], qty, price)
	case "list-orders":
		listOrders()
	case "create-order":
		if len(os.Args) < 3 {
			fmt.Println("usage: create-order ITEM_ID:QTY[,ITEM_ID:QTY]")
			os.Exit(1)
		}
		createOrder(os.Args[2])
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("client commands:")
	fmt.Println("  list-items")
	fmt.Println("  create-item NAME QUANTITY PRICE")
	fmt.Println("  list-orders")
	fmt.Println("  create-order ITEM_ID:QTY[,ITEM_ID:QTY]")
}

func listItems() {
	resp, err := http.Get("http://localhost:8001/items")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func createItem(name string, qty int, price float64) {
	b, _ := json.Marshal(map[string]interface{}{"name": name, "quantity": qty, "price": price})
	resp, err := http.Post("http://localhost:8001/items", "application/json", bytes.NewReader(b))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func listOrders() {
	resp, err := http.Get("http://localhost:8002/orders")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func createOrder(spec string) {
	parts := strings.Split(spec, ",")
	items := make([]map[string]int, 0, len(parts))
	for _, p := range parts {
		kv := strings.Split(p, ":")
		if len(kv) != 2 {
			fmt.Println("invalid item spec:", p)
			return
		}
		id, _ := strconv.Atoi(kv[0])
		q, _ := strconv.Atoi(kv[1])
		items = append(items, map[string]int{"id": id, "quantity": q})
	}
	b, _ := json.Marshal(map[string]interface{}{"items": items})
	resp, err := http.Post("http://localhost:8002/orders", "application/json", bytes.NewReader(b))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
