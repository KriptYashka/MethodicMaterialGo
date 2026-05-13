package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
)

type Item struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

var items = []Item{
    {ID: 1, Name: "Item 1"},
    {ID: 2, Name: "Item 2"},
}
var nextID = 3

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func handleItems(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        writeJSON(w, http.StatusOK, items)
    case http.MethodPost:
        var item Item
        if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
            return
        }
        item.ID = nextID
        nextID++
        items = append(items, item)
        writeJSON(w, http.StatusCreated, item)
    default:
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
    }
}

func handleItem(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
        return
    }

    var found *Item
    for i, item := range items {
        if item.ID == id {
            found = &items[i]
            break
        }
    }
    if found == nil {
        writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
        return
    }

    switch r.Method {
    case http.MethodGet:
        writeJSON(w, http.StatusOK, found)
    case http.MethodDelete:
        idx := -1
        for i, item := range items {
            if item.ID == id {
                idx = i
                break
            }
        }
        items = append(items[:idx], items[idx+1:]...)
        writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
    default:
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
    }
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /items", handleItems)
    mux.HandleFunc("POST /items", handleItems)
    mux.HandleFunc("GET /items/{id}", handleItem)
    mux.HandleFunc("DELETE /items/{id}", handleItem)

    addr := ":8080"
    fmt.Printf("Router demo on %s\n", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
}

// Test:
//   curl http://localhost:8080/items
//   curl -X POST -d '{"name":"Item 3"}' http://localhost:8080/items
//   curl http://localhost:8080/items/1
//   curl -X DELETE http://localhost:8080/items/1
